package spotify

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"audioinsights-service-go/internal/domain"
)

func TestFetchProfileForwardsAuthorization(t *testing.T) {
	var gotAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/me" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		gotAuth = r.Header.Get("Authorization")
		writeFixture(t, w, domain.UserProfile{DisplayName: "John Doe", Images: []domain.Image{}})
	}))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	profile, err := client.FetchProfile(context.Background(), "Bearer mockToken")
	if err != nil {
		t.Fatalf("FetchProfile returned error: %v", err)
	}
	if gotAuth != "Bearer mockToken" {
		t.Fatalf("expected auth forwarded, got %q", gotAuth)
	}
	if profile.DisplayName != "John Doe" {
		t.Fatalf("expected decoded profile, got %#v", profile)
	}
}

func TestFetchTopItemsForwardsTimeRange(t *testing.T) {
	seen := map[string]string{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer mockToken" {
			t.Fatalf("expected auth forwarded, got %q", r.Header.Get("Authorization"))
		}
		seen[r.URL.Path] = r.URL.Query().Get("time_range")

		switch r.URL.Path {
		case "/v1/me/top/tracks":
			writeFixture(t, w, domain.TopItemsResponse[domain.Track]{Items: []domain.Track{{ID: "track-1", Name: "Track"}}})
		case "/v1/me/top/artists":
			writeFixture(t, w, domain.TopItemsResponse[domain.Artist]{Items: []domain.Artist{{ID: "artist-1", Name: "Artist", Genres: []string{"pop", "rock"}}}})
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	topItems, err := client.FetchTopItems(context.Background(), "Bearer mockToken", "short_term")
	if err != nil {
		t.Fatalf("FetchTopItems returned error: %v", err)
	}

	if seen["/v1/me/top/tracks"] != "short_term" {
		t.Fatalf("expected tracks time_range short_term, got %q", seen["/v1/me/top/tracks"])
	}
	if seen["/v1/me/top/artists"] != "short_term" {
		t.Fatalf("expected artists time_range short_term, got %q", seen["/v1/me/top/artists"])
	}
	if len(topItems.Tracks) != 1 || len(topItems.Artists) != 1 {
		t.Fatalf("expected combined top items, got %#v", topItems)
	}
	if len(topItems.Genres) != 2 || topItems.Genres[0] != "pop" || topItems.Genres[1] != "rock" {
		t.Fatalf("expected flattened genres, got %#v", topItems.Genres)
	}
}

func TestSpotifyNonSuccessStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "nope", http.StatusUnauthorized)
	}))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	_, err := client.FetchProfile(context.Background(), "Bearer bad")
	if err == nil {
		t.Fatal("expected error")
	}

	spotifyErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("expected spotify error, got %T", err)
	}
	if spotifyErr.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, spotifyErr.StatusCode)
	}
}

func writeFixture(t *testing.T, w http.ResponseWriter, value any) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(value); err != nil {
		t.Fatalf("encode fixture: %v", err)
	}
}
