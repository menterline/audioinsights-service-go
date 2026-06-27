package domain

import (
	"encoding/json"
	"os"
	"testing"
)

func TestUserProfileFixtureCompatibility(t *testing.T) {
	data, err := os.ReadFile("../../testdata/dummyUserProfile.json")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	var profile UserProfile
	if err := json.Unmarshal(data, &profile); err != nil {
		t.Fatalf("unmarshal profile: %v", err)
	}

	if profile.DisplayName == "" || profile.ExternalURLs.Spotify == "" {
		t.Fatalf("expected populated profile, got %#v", profile)
	}
}

func TestTopItemsFixtureCompatibility(t *testing.T) {
	tracksData, err := os.ReadFile("../../testdata/dummyTopTracksResponse.json")
	if err != nil {
		t.Fatalf("read tracks fixture: %v", err)
	}
	artistsData, err := os.ReadFile("../../testdata/dummyTopArtistsResponse.json")
	if err != nil {
		t.Fatalf("read artists fixture: %v", err)
	}

	var tracks TopItemsResponse[Track]
	if err := json.Unmarshal(tracksData, &tracks); err != nil {
		t.Fatalf("unmarshal tracks: %v", err)
	}

	var artists TopItemsResponse[Artist]
	if err := json.Unmarshal(artistsData, &artists); err != nil {
		t.Fatalf("unmarshal artists: %v", err)
	}

	topItems := NewTopItems("short_term", artists, tracks)
	if topItems.Term != "short_term" {
		t.Fatalf("expected term, got %q", topItems.Term)
	}
	if len(topItems.Tracks) == 0 || len(topItems.Artists) == 0 || len(topItems.Genres) == 0 {
		t.Fatalf("expected populated top items, got %#v", topItems)
	}
}
