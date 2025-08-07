package musicbrainz

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"relay/app/cache"
)

const (
	baseURL   = "https://musicbrainz.org/ws/2"
	userAgent = "MediaManager-MetadataRelay/1.0 (https://github.com/np22-jpg/MediaManager)"
)

// initializes the MusicBrainz client
func InitMusicBrainz() {
	fmt.Printf("MusicBrainz client initialized\n")
}

// makes an HTTP request to MusicBrainz API
func makeRequest(endpoint string, params url.Values) (any, error) {
	// Add required format parameter
	params.Set("fmt", "json")

	url := fmt.Sprintf("%s%s?%s", baseURL, endpoint, params.Encode())

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// MusicBrainz requires a User-Agent header
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
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

// searches for artists
func SearchArtists(ctx context.Context, query string, limit int) (any, error) {
	return cache.NewCache("musicbrainz_artist_search").TTL(24*time.Hour).Wrap(func() (any, error) {
		params := url.Values{}
		params.Set("query", query)
		params.Set("limit", fmt.Sprintf("%d", limit))
		return makeRequest("/artist", params)
	})(ctx, query, limit)
}

// gets a specific artist by MBID
func GetArtist(ctx context.Context, mbid string) (any, error) {
	return cache.NewCache("musicbrainz_artist").TTL(7*24*time.Hour).Wrap(func() (any, error) {
		params := url.Values{}
		params.Set("inc", "aliases+tags+ratings+genres")
		return makeRequest(fmt.Sprintf("/artist/%s", mbid), params)
	})(ctx, mbid)
}

// searches for release groups (albums)
func SearchReleaseGroups(ctx context.Context, query string, limit int) (any, error) {
	return cache.NewCache("musicbrainz_release_group_search").TTL(24*time.Hour).Wrap(func() (any, error) {
		params := url.Values{}
		params.Set("query", query)
		params.Set("limit", fmt.Sprintf("%d", limit))
		return makeRequest("/release-group", params)
	})(ctx, query, limit)
}

// gets a specific release group by MBID
func GetReleaseGroup(ctx context.Context, mbid string) (any, error) {
	return cache.NewCache("musicbrainz_release_group").TTL(7*24*time.Hour).Wrap(func() (any, error) {
		params := url.Values{}
		params.Set("inc", "aliases+tags+ratings+genres+releases")
		return makeRequest(fmt.Sprintf("/release-group/%s", mbid), params)
	})(ctx, mbid)
}

// searches for releases
func SearchReleases(ctx context.Context, query string, limit int) (any, error) {
	return cache.NewCache("musicbrainz_release_search").TTL(24*time.Hour).Wrap(func() (any, error) {
		params := url.Values{}
		params.Set("query", query)
		params.Set("limit", fmt.Sprintf("%d", limit))
		return makeRequest("/release", params)
	})(ctx, query, limit)
}

// gets a specific release by MBID
func GetRelease(ctx context.Context, mbid string) (any, error) {
	return cache.NewCache("musicbrainz_release").TTL(7*24*time.Hour).Wrap(func() (any, error) {
		params := url.Values{}
		params.Set("inc", "aliases+tags+ratings+genres+recordings+artist-credits")
		return makeRequest(fmt.Sprintf("/release/%s", mbid), params)
	})(ctx, mbid)
}

// searches for recordings (tracks)
func SearchRecordings(ctx context.Context, query string, limit int) (any, error) {
	return cache.NewCache("musicbrainz_recording_search").TTL(24*time.Hour).Wrap(func() (any, error) {
		params := url.Values{}
		params.Set("query", query)
		params.Set("limit", fmt.Sprintf("%d", limit))
		return makeRequest("/recording", params)
	})(ctx, query, limit)
}

// gets a specific recording by MBID
func GetRecording(ctx context.Context, mbid string) (any, error) {
	return cache.NewCache("musicbrainz_recording").TTL(7*24*time.Hour).Wrap(func() (any, error) {
		params := url.Values{}
		params.Set("inc", "aliases+tags+ratings+genres+releases+artist-credits")
		return makeRequest(fmt.Sprintf("/recording/%s", mbid), params)
	})(ctx, mbid)
}

// browses release groups for a specific artist
func BrowseArtistReleaseGroups(ctx context.Context, artistMbid string, limit int) (any, error) {
	return cache.NewCache("musicbrainz_artist_release_groups").TTL(24*time.Hour).Wrap(func() (any, error) {
		params := url.Values{}
		params.Set("artist", artistMbid)
		params.Set("limit", fmt.Sprintf("%d", limit))
		params.Set("inc", "tags+ratings+genres")
		return makeRequest("/release-group", params)
	})(ctx, artistMbid, limit)
}

// browses releases for a specific release group
func BrowseReleaseGroupReleases(ctx context.Context, releaseGroupMbid string, limit int) (any, error) {
	return cache.NewCache("musicbrainz_release_group_releases").TTL(24*time.Hour).Wrap(func() (any, error) {
		params := url.Values{}
		params.Set("release-group", releaseGroupMbid)
		params.Set("limit", fmt.Sprintf("%d", limit))
		params.Set("inc", "tags+ratings+genres")
		return makeRequest("/release", params)
	})(ctx, releaseGroupMbid, limit)
}

// performs an advanced artist search with field-specific queries
func AdvancedSearchArtists(ctx context.Context, artistName, area, beginDate, endDate string, limit int) (any, error) {
	cacheKey := fmt.Sprintf("musicbrainz_artist_advanced_search_%s_%s_%s_%s_%d",
		artistName, area, beginDate, endDate, limit)

	return cache.NewCache(cacheKey).TTL(24*time.Hour).Wrap(func() (any, error) {
		var queryParts []string

		if artistName != "" {
			queryParts = append(queryParts, fmt.Sprintf("artist:\"%s\"", artistName))
		}
		if area != "" {
			queryParts = append(queryParts, fmt.Sprintf("area:\"%s\"", area))
		}
		if beginDate != "" {
			queryParts = append(queryParts, fmt.Sprintf("begin:%s", beginDate))
		}
		if endDate != "" {
			queryParts = append(queryParts, fmt.Sprintf("end:%s", endDate))
		}

		query := strings.Join(queryParts, " AND ")

		params := url.Values{}
		params.Set("query", query)
		params.Set("limit", fmt.Sprintf("%d", limit))
		return makeRequest("/artist", params)
	})(ctx, artistName, area, beginDate, endDate, limit)
}
