package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"audioinsights-service-go/internal/domain"
)

const BaseURL = "https://api.spotify.com"

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type Error struct {
	StatusCode int
	Message    string
}

func (e *Error) Error() string {
	return e.Message
}

func NewClient() *Client {
	return NewClientWithBaseURL(BaseURL)
}

func NewClientWithBaseURL(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) FetchProfile(ctx context.Context, bearerToken string) (domain.UserProfile, error) {
	var profile domain.UserProfile
	err := c.getJSON(ctx, "/v1/me", nil, bearerToken, &profile)
	return profile, err
}

func (c *Client) FetchTopTracks(ctx context.Context, bearerToken, term string) (domain.TopItemsResponse[domain.Track], error) {
	var tracks domain.TopItemsResponse[domain.Track]
	err := c.getJSON(ctx, "/v1/me/top/tracks", url.Values{"time_range": []string{term}}, bearerToken, &tracks)
	return tracks, err
}

func (c *Client) FetchTopArtists(ctx context.Context, bearerToken, term string) (domain.TopItemsResponse[domain.Artist], error) {
	var artists domain.TopItemsResponse[domain.Artist]
	err := c.getJSON(ctx, "/v1/me/top/artists", url.Values{"time_range": []string{term}}, bearerToken, &artists)
	return artists, err
}

func (c *Client) FetchTopItems(ctx context.Context, bearerToken, term string) (domain.TopItems, error) {
	tracks, err := c.FetchTopTracks(ctx, bearerToken, term)
	if err != nil {
		return domain.TopItems{}, err
	}

	artists, err := c.FetchTopArtists(ctx, bearerToken, term)
	if err != nil {
		return domain.TopItems{}, err
	}

	return domain.NewTopItems(term, artists, tracks), nil
}

func (c *Client) getJSON(ctx context.Context, path string, query url.Values, bearerToken string, target any) error {
	reqURL := c.baseURL + path
	if len(query) > 0 {
		reqURL += "?" + query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return fmt.Errorf("create spotify request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("call spotify: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return &Error{StatusCode: resp.StatusCode, Message: fmt.Sprintf("spotify returned status %d", resp.StatusCode)}
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("decode spotify response: %w", err)
	}

	return nil
}
