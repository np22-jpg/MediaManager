package seadex

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"relay/app/cache"
)

var (
	baseURL string
)

// InitSeaDex initializes the SeaDex client with the provided base URL.
func InitSeaDex(url string) {
	baseURL = url
	if baseURL == "" {
		baseURL = "https://releases.moe"
	}
	slog.Info("SeaDex initialized", "baseUrl", baseURL)
}

// EntryRecord represents a single anime entry in SeaDex.
type EntryRecord struct {
	ID              string          `json:"id"`
	AnilistID       int             `json:"anilist_id"`
	CollectionID    string          `json:"collection_id"`
	CollectionName  string          `json:"collection_name"`
	Comparisons     []string        `json:"comparisons"`
	CreatedAt       string          `json:"created"`
	IsIncomplete    bool            `json:"is_incomplete"`
	Notes           string          `json:"notes"`
	Size            int64           `json:"size"`
	TheoreticalBest *string         `json:"theoretical_best"`
	Torrents        []TorrentRecord `json:"torrents"`
	UpdatedAt       string          `json:"updated"`
	URL             string          `json:"url"`
}

// TorrentRecord represents a single torrent record within a SeaDex entry.
type TorrentRecord struct {
	ID             string   `json:"id"`
	CollectionID   string   `json:"collection_id"`
	CollectionName string   `json:"collection_name"`
	CreatedAt      string   `json:"created"`
	Files          []File   `json:"files"`
	GroupedURL     *string  `json:"grouped_url"`
	Infohash       *string  `json:"infohash"`
	IsBest         bool     `json:"is_best"`
	IsDualAudio    bool     `json:"is_dual_audio"`
	ReleaseGroup   string   `json:"release_group"`
	Size           int64    `json:"size"`
	Tags           []string `json:"tags"`
	Tracker        string   `json:"tracker"`
	UpdatedAt      string   `json:"updated"`
	URL            string   `json:"url"`
}

// File represents a file in the torrent.
type File struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

// makeRequest makes an HTTP request to SeaDex API with proper error handling.
func makeRequest(endpoint string, params url.Values) (any, error) {
	reqURL := fmt.Sprintf("%s/api/collections/entries/records%s", baseURL, endpoint)
	if len(params) > 0 {
		reqURL += "?" + params.Encode()
	}

	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			slog.Error("failed to close response body", "error", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var result any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// SearchEntries searches for anime entries with query, cached for 4 hours.
func SearchEntries(ctx context.Context, query string, page int, perPage int) (any, error) {
	return cache.NewCache("seadx_search_entries").TTL(4*time.Hour).Wrap(func() (any, error) {
		params := url.Values{}
		if query != "" {
			// SeaDex uses PocketBase-style filtering
			params.Set("filter", fmt.Sprintf("(collection_name ~ '%s' || torrents.release_group ~ '%s')", query, query))
		}
		if page > 0 {
			params.Set("page", fmt.Sprintf("%d", page))
		}
		if perPage > 0 {
			params.Set("perPage", fmt.Sprintf("%d", perPage))
		}
		return makeRequest("", params)
	})(ctx, query, page, perPage)
}

// GetEntryByID retrieves an entry by its ID with 8-hour caching.
func GetEntryByID(ctx context.Context, id string) (any, error) {
	return cache.NewCache("seadx_entry_by_id").TTL(8*time.Hour).Wrap(func() (any, error) {
		return makeRequest("/"+id, nil)
	})(ctx, id)
}

// GetEntryByAnilistID retrieves an entry by AniList ID with 8-hour caching.
func GetEntryByAnilistID(ctx context.Context, anilistID int) (any, error) {
	return cache.NewCache("seadx_entry_by_anilist").TTL(8*time.Hour).Wrap(func() (any, error) {
		params := url.Values{}
		params.Set("filter", fmt.Sprintf("anilist_id = %d", anilistID))
		return makeRequest("", params)
	})(ctx, anilistID)
}

// GetTrendingEntries gets trending/popular anime entries, cached for 2 hours.
func GetTrendingEntries(ctx context.Context, limit int) (any, error) {
	return cache.NewCache("seadx_trending").TTL(2*time.Hour).Wrap(func() (any, error) {
		params := url.Values{}
		params.Set("sort", "-updated") // Sort by recently updated
		if limit > 0 {
			params.Set("perPage", fmt.Sprintf("%d", limit))
		}
		return makeRequest("", params)
	})(ctx, limit)
}

// GetEntriesByReleaseGroup searches for entries by release group with 6-hour caching.
func GetEntriesByReleaseGroup(ctx context.Context, releaseGroup string) (any, error) {
	return cache.NewCache("seadx_by_release_group").TTL(6*time.Hour).Wrap(func() (any, error) {
		params := url.Values{}
		params.Set("filter", fmt.Sprintf("torrents.release_group ~ '%s'", releaseGroup))
		return makeRequest("", params)
	})(ctx, releaseGroup)
}

// GetEntriesByTracker searches for entries by tracker with 6-hour caching.
func GetEntriesByTracker(ctx context.Context, tracker string) (any, error) {
	return cache.NewCache("seadx_by_tracker").TTL(6*time.Hour).Wrap(func() (any, error) {
		params := url.Values{}
		params.Set("filter", fmt.Sprintf("torrents.tracker ~ '%s'", tracker))
		return makeRequest("", params)
	})(ctx, tracker)
}
