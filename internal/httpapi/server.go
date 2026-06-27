package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"audioinsights-service-go/internal/domain"
	"audioinsights-service-go/internal/spotify"
)

type SpotifyService interface {
	FetchProfile(ctx context.Context, bearerToken string) (domain.UserProfile, error)
	FetchTopItems(ctx context.Context, bearerToken, term string) (domain.TopItems, error)
}

type Server struct {
	spotify SpotifyService
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewServer(spotify SpotifyService) *Server {
	return &Server{spotify: spotify}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	addCORSHeaders(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	switch {
	case r.Method == http.MethodGet && r.URL.Path == "/api/profile/":
		s.fetchProfile(w, r)
	case r.Method == http.MethodGet && r.URL.Path == "/api/profile/topItems":
		s.fetchTopItems(w, r)
	default:
		writeError(w, http.StatusNotFound, "not found")
	}
}

func (s *Server) fetchProfile(w http.ResponseWriter, r *http.Request) {
	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		writeError(w, http.StatusUnauthorized, "missing Authorization header")
		return
	}

	profile, err := s.spotify.FetchProfile(r.Context(), bearerToken)
	if err != nil {
		writeUpstreamError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, profile)
}

func (s *Server) fetchTopItems(w http.ResponseWriter, r *http.Request) {
	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		writeError(w, http.StatusUnauthorized, "missing Authorization header")
		return
	}

	term := r.URL.Query().Get("term")
	if term == "" {
		writeError(w, http.StatusBadRequest, "missing term query parameter")
		return
	}

	topItems, err := s.spotify.FetchTopItems(r.Context(), bearerToken, term)
	if err != nil {
		writeUpstreamError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, topItems)
}

func addCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
}

func writeUpstreamError(w http.ResponseWriter, err error) {
	var spotifyErr *spotify.Error
	if errors.As(err, &spotifyErr) && spotifyErr.StatusCode >= 400 && spotifyErr.StatusCode <= 599 {
		writeError(w, spotifyErr.StatusCode, spotifyErr.Message)
		return
	}

	writeError(w, http.StatusBadGateway, "spotify request failed")
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}
