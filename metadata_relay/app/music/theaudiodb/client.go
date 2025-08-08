package theaudiodb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"relay/app/cache"
)

type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	BaseURL   string
	APIKey    string
	HTTP      HTTPDoer
	Timeout   time.Duration
	UserAgent string
}

func New(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL:   baseURL,
		APIKey:    apiKey,
		HTTP:      &http.Client{Timeout: 8 * time.Second},
		Timeout:   8 * time.Second,
		UserAgent: "metadata-relay/1.0",
	}
}

// Minimal response shapes based on TheAudioDB docs
type artistSearchResponse struct {
	Artists []struct {
		IDArtist         string `json:"idArtist"`
		StrArtist        string `json:"strArtist"`
		StrBiographyEN   string `json:"strBiographyEN"`
		StrWebsite       string `json:"strWebsite"`
		StrGenre         string `json:"strGenre"`
		IntFormedYear    string `json:"intFormedYear"`
		StrMusicBrainzID string `json:"strMusicBrainzID"`
		StrArtistThumb   string `json:"strArtistThumb"`
		StrArtistFanart  string `json:"strArtistFanart"`
		StrCountry       string `json:"strCountry"`
	} `json:"artists"`
}

type albumSearchResponse struct {
	Album []struct {
		IDAlbum          string `json:"idAlbum"`
		IDMBAlbum        string `json:"idMBAlbum"`
		StrAlbum         string `json:"strAlbum"`
		StrArtist        string `json:"strArtist"`
		IntYearReleased  string `json:"intYearReleased"`
		StrGenre         string `json:"strGenre"`
		StrAlbumThumb    string `json:"strAlbumThumb"`
		StrDescriptionEN string `json:"strDescriptionEN"`
	} `json:"album"`
}

type trackSearchResponse struct {
	Track []struct {
		IDTrack          string `json:"idTrack"`
		IDMBTrack        string `json:"idMBTrack"`
		StrTrack         string `json:"strTrack"`
		StrArtist        string `json:"strArtist"`
		StrAlbum         string `json:"strAlbum"`
		IntDuration      string `json:"intDuration"`
		StrGenre         string `json:"strGenre"`
		StrDescriptionEN string `json:"strDescriptionEN"`
		StrTrackThumb    string `json:"strTrackThumb"`
	} `json:"track"`
}

// SearchArtist returns best match by name
func (c *Client) SearchArtist(ctx context.Context, name string) (map[string]any, error) {
	if c == nil || c.BaseURL == "" || c.APIKey == "" {
		return nil, fmt.Errorf("theaudiodb not configured")
	}
	base := fmt.Sprintf("%s/%s/search.php", c.BaseURL, c.APIKey)
	u, err := url.Parse(base)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("s", name)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("theaudiodb http %d", resp.StatusCode)
	}

	var out artistSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	if len(out.Artists) == 0 {
		return map[string]any{"count": 0, "artists": []any{}}, nil
	}
	// Return a simplified first result
	a := out.Artists[0]
	return map[string]any{
		"name":       a.StrArtist,
		"biography":  a.StrBiographyEN,
		"website":    a.StrWebsite,
		"genre":      a.StrGenre,
		"formedYear": a.IntFormedYear,
		"mbid":       a.StrMusicBrainzID,
		"thumb":      a.StrArtistThumb,
		"fanart":     a.StrArtistFanart,
		"country":    a.StrCountry,
	}, nil
}

// SearchArtistByMBID returns artist data by MusicBrainz ID
func (c *Client) SearchArtistByMBID(ctx context.Context, mbid string) (map[string]any, error) {
	if c == nil || c.BaseURL == "" || c.APIKey == "" {
		return nil, fmt.Errorf("theaudiodb not configured")
	}
	base := fmt.Sprintf("%s/%s/artist-mb.php", c.BaseURL, c.APIKey)
	u, err := url.Parse(base)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("i", mbid)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("theaudiodb http %d", resp.StatusCode)
	}

	var out artistSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	if len(out.Artists) == 0 {
		return map[string]any{"error": "artist not found"}, nil
	}
	// Return the first result
	a := out.Artists[0]
	return map[string]any{
		"name":       a.StrArtist,
		"biography":  a.StrBiographyEN,
		"website":    a.StrWebsite,
		"genre":      a.StrGenre,
		"formedYear": a.IntFormedYear,
		"mbid":       a.StrMusicBrainzID,
		"thumb":      a.StrArtistThumb,
		"fanart":     a.StrArtistFanart,
		"country":    a.StrCountry,
	}, nil
}

// SearchAlbumByMBID returns album data by MusicBrainz release group ID
func (c *Client) SearchAlbumByMBID(ctx context.Context, mbid string) (map[string]any, error) {
	if c == nil || c.BaseURL == "" || c.APIKey == "" {
		return nil, fmt.Errorf("theaudiodb not configured")
	}
	base := fmt.Sprintf("%s/%s/album-mb.php", c.BaseURL, c.APIKey)
	u, err := url.Parse(base)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("i", mbid)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("theaudiodb http %d", resp.StatusCode)
	}

	var out albumSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	if len(out.Album) == 0 {
		return map[string]any{"error": "album not found"}, nil
	}
	// Return the first result
	a := out.Album[0]
	return map[string]any{
		"name":        a.StrAlbum,
		"artist":      a.StrArtist,
		"year":        a.IntYearReleased,
		"genre":       a.StrGenre,
		"thumb":       a.StrAlbumThumb,
		"description": a.StrDescriptionEN,
		"mbid":        a.IDMBAlbum,
	}, nil
}

// SearchTrackByMBID returns track data by MusicBrainz recording ID
func (c *Client) SearchTrackByMBID(ctx context.Context, mbid string) (map[string]any, error) {
	if c == nil || c.BaseURL == "" || c.APIKey == "" {
		return nil, fmt.Errorf("theaudiodb not configured")
	}
	base := fmt.Sprintf("%s/%s/track-mb.php", c.BaseURL, c.APIKey)
	u, err := url.Parse(base)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("i", mbid)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("theaudiodb http %d", resp.StatusCode)
	}

	var out trackSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	if len(out.Track) == 0 {
		return map[string]any{"error": "track not found"}, nil
	}
	// Return the first result
	t := out.Track[0]
	return map[string]any{
		"name":        t.StrTrack,
		"artist":      t.StrArtist,
		"album":       t.StrAlbum,
		"duration":    t.IntDuration,
		"genre":       t.StrGenre,
		"description": t.StrDescriptionEN,
		"thumb":       t.StrTrackThumb,
		"mbid":        t.IDMBTrack,
	}, nil
}

// Cached wrapper functions with prioritized caching for MBID lookups

// GetArtistByMBID gets artist by MBID with high-priority caching (7 day TTL)
func GetArtistByMBID(ctx context.Context, mbid string) (any, error) {
	return cache.NewCache("theaudiodb:artist_mbid").TTL(7*24*time.Hour).Wrap(func() (any, error) {
		return client.SearchArtistByMBID(ctx, mbid)
	})(ctx, mbid)
}

// GetAlbumByMBID gets album by MBID with high-priority caching (7 day TTL)
func GetAlbumByMBID(ctx context.Context, mbid string) (any, error) {
	return cache.NewCache("theaudiodb:album_mbid").TTL(7*24*time.Hour).Wrap(func() (any, error) {
		return client.SearchAlbumByMBID(ctx, mbid)
	})(ctx, mbid)
}

// GetTrackByMBID gets track by MBID with high-priority caching (7 day TTL)
func GetTrackByMBID(ctx context.Context, mbid string) (any, error) {
	return cache.NewCache("theaudiodb:track_mbid").TTL(7*24*time.Hour).Wrap(func() (any, error) {
		return client.SearchTrackByMBID(ctx, mbid)
	})(ctx, mbid)
}

// SearchArtist gets artist by name with standard caching (24 hour TTL)
func SearchArtist(ctx context.Context, name string) (any, error) {
	return cache.NewCache("theaudiodb:artist_search").TTL(24*time.Hour).Wrap(func() (any, error) {
		return client.SearchArtist(ctx, name)
	})(ctx, name)
}
