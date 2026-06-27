package domain

type ExplicitContent struct {
	FilterEnabled bool `json:"filter_enabled"`
	FilterLocked  bool `json:"filter_locked"`
}

type ExternalURL struct {
	Spotify string `json:"spotify"`
}

type Followers struct {
	Href  *string `json:"href"`
	Total int     `json:"total"`
}

type Image struct {
	URL    string `json:"url"`
	Height *int   `json:"height"`
	Width  *int   `json:"width"`
}

type UserProfile struct {
	Country         string          `json:"country"`
	DisplayName     string          `json:"display_name"`
	Email           string          `json:"email"`
	ExplicitContent ExplicitContent `json:"explicit_content"`
	ExternalURLs    ExternalURL     `json:"external_urls"`
	Followers       Followers       `json:"followers"`
	Href            string          `json:"href"`
	ID              string          `json:"id"`
	Images          []Image         `json:"images"`
	Product         string          `json:"product"`
	Type            string          `json:"type"`
	URI             string          `json:"uri"`
}

type Artist struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Genres       []string    `json:"genres"`
	Href         string      `json:"href"`
	ExternalURLs ExternalURL `json:"external_urls"`
	Images       []Image     `json:"images"`
	Popularity   int         `json:"popularity"`
	Type         string      `json:"type"`
	URI          string      `json:"uri"`
}

type Album struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Images       []Image     `json:"images"`
	ReleaseDate  string      `json:"release_date"`
	TotalTracks  int         `json:"total_tracks"`
	Href         string      `json:"href"`
	ExternalURLs ExternalURL `json:"external_urls"`
	URI          string      `json:"uri"`
}

type Track struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Album        Album       `json:"album"`
	Artists      []Artist    `json:"artists"`
	DurationMS   int         `json:"duration_ms"`
	Explicit     bool        `json:"explicit"`
	Href         string      `json:"href"`
	ExternalURLs ExternalURL `json:"external_urls"`
	Popularity   int         `json:"popularity"`
	PreviewURL   *string     `json:"preview_url"`
	TrackNumber  int         `json:"track_number"`
	Type         string      `json:"type"`
	URI          string      `json:"uri"`
}

type TopItemsResponse[T any] struct {
	Items    []T     `json:"items"`
	Total    int     `json:"total"`
	Limit    int     `json:"limit"`
	Offset   int     `json:"offset"`
	Href     string  `json:"href"`
	Previous *string `json:"previous"`
	Next     *string `json:"next"`
}

type TopItems struct {
	Term    string   `json:"term"`
	Artists []Artist `json:"artists"`
	Tracks  []Track  `json:"tracks"`
	Genres  []string `json:"genres"`
}

func NewTopItems(term string, artists TopItemsResponse[Artist], tracks TopItemsResponse[Track]) TopItems {
	genres := make([]string, 0)
	for _, artist := range artists.Items {
		genres = append(genres, artist.Genres...)
	}

	return TopItems{
		Term:    term,
		Artists: artists.Items,
		Tracks:  tracks.Items,
		Genres:  genres,
	}
}
