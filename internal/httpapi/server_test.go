package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"audioinsights-service-go/internal/domain"
	"audioinsights-service-go/internal/spotify"
)

type stubSpotify struct {
	profileAuth  string
	topItemsAuth string
	topItemsTerm string
	profile      domain.UserProfile
	topItems     domain.TopItems
	err          error
}

func (s *stubSpotify) FetchProfile(_ context.Context, bearerToken string) (domain.UserProfile, error) {
	s.profileAuth = bearerToken
	return s.profile, s.err
}

func (s *stubSpotify) FetchTopItems(_ context.Context, bearerToken, term string) (domain.TopItems, error) {
	s.topItemsAuth = bearerToken
	s.topItemsTerm = term
	return s.topItems, s.err
}

func TestFetchProfile(t *testing.T) {
	service := &stubSpotify{
		profile: domain.UserProfile{
			Country:     "US",
			DisplayName: "John Doe",
			Email:       "john@example.com",
			ID:          "john",
			Images:      []domain.Image{},
		},
	}
	server := NewServer(service)

	req := httptest.NewRequest(http.MethodGet, "/api/profile/", nil)
	req.Header.Set("Authorization", "Bearer mockToken")
	rr := httptest.NewRecorder()

	server.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	if service.profileAuth != "Bearer mockToken" {
		t.Fatalf("expected forwarded auth, got %q", service.profileAuth)
	}

	var got domain.UserProfile
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.DisplayName != "John Doe" {
		t.Fatalf("expected profile response, got %#v", got)
	}
}

func TestFetchTopItems(t *testing.T) {
	service := &stubSpotify{
		topItems: domain.TopItems{
			Term:   "short_term",
			Genres: []string{"pop", "rock"},
		},
	}
	server := NewServer(service)

	req := httptest.NewRequest(http.MethodGet, "/api/profile/topItems?term=short_term", nil)
	req.Header.Set("Authorization", "Bearer mockToken")
	rr := httptest.NewRecorder()

	server.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	if service.topItemsAuth != "Bearer mockToken" {
		t.Fatalf("expected forwarded auth, got %q", service.topItemsAuth)
	}
	if service.topItemsTerm != "short_term" {
		t.Fatalf("expected term short_term, got %q", service.topItemsTerm)
	}
}

func TestValidationErrors(t *testing.T) {
	server := NewServer(&stubSpotify{})

	tests := []struct {
		name   string
		path   string
		auth   string
		status int
	}{
		{name: "missing auth", path: "/api/profile/", status: http.StatusUnauthorized},
		{name: "missing term", path: "/api/profile/topItems", auth: "Bearer token", status: http.StatusBadRequest},
		{name: "tracks analysis removed", path: "/api/profile/tracksAnalysis?ids=1", auth: "Bearer token", status: http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			if tt.auth != "" {
				req.Header.Set("Authorization", tt.auth)
			}
			rr := httptest.NewRecorder()

			server.ServeHTTP(rr, req)

			if rr.Code != tt.status {
				t.Fatalf("expected status %d, got %d", tt.status, rr.Code)
			}
		})
	}
}

func TestSpotifyErrors(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		status int
	}{
		{name: "spotify status proxied", err: &spotify.Error{StatusCode: http.StatusForbidden, Message: "spotify returned status 403"}, status: http.StatusForbidden},
		{name: "generic upstream failure", err: errors.New("boom"), status: http.StatusBadGateway},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewServer(&stubSpotify{err: tt.err})
			req := httptest.NewRequest(http.MethodGet, "/api/profile/", nil)
			req.Header.Set("Authorization", "Bearer token")
			rr := httptest.NewRecorder()

			server.ServeHTTP(rr, req)

			if rr.Code != tt.status {
				t.Fatalf("expected status %d, got %d", tt.status, rr.Code)
			}
		})
	}
}
