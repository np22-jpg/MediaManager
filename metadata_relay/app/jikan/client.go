package jikan

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"relay/app/cache"
)

var (
	baseURL string
)

// InitJikan initializes the Jikan client with the provided base URL.
func InitJikan(url string) {
	baseURL = url
	if baseURL == "" {
		baseURL = "https://api.jikan.moe/v4"
	}
	slog.Info("Jikan initialized", "baseUrl", baseURL)
}

// AnimeData represents anime information from Jikan API.
type AnimeData struct {
	ID       int    `json:"mal_id"`
	URL      string `json:"url"`
	Title    string `json:"title"`
	TitleJP  string `json:"title_japanese"`
	TitleEN  string `json:"title_english"`
	Type     string `json:"type"`
	Episodes int    `json:"episodes"`
	Status   string `json:"status"`
	Aired    struct {
		From string `json:"from"`
		To   string `json:"to"`
	} `json:"aired"`
	Duration string  `json:"duration"`
	Rating   string  `json:"rating"`
	Score    float64 `json:"score"`
	Synopsis string  `json:"synopsis"`
	Images   struct {
		JPG struct {
			ImageURL      string `json:"image_url"`
			SmallImageURL string `json:"small_image_url"`
			LargeImageURL string `json:"large_image_url"`
		} `json:"jpg"`
	} `json:"images"`
	Genres []struct {
		ID   int    `json:"mal_id"`
		Name string `json:"name"`
	} `json:"genres"`
	Studios []struct {
		ID   int    `json:"mal_id"`
		Name string `json:"name"`
	} `json:"studios"`
}

// SearchResponse represents search results from Jikan API.
type SearchResponse struct {
	Data       []AnimeData `json:"data"`
	Pagination struct {
		LastVisiblePage int  `json:"last_visible_page"`
		HasNextPage     bool `json:"has_next_page"`
		CurrentPage     int  `json:"current_page"`
	} `json:"pagination"`
}

// SingleResponse represents a single anime response.
type SingleResponse struct {
	Data AnimeData `json:"data"`
}

// RecommendationsResponse represents anime recommendations.
type RecommendationsResponse struct {
	Data []struct {
		Entry AnimeData `json:"entry"`
	} `json:"data"`
}

// makeRequest makes an HTTP request to Jikan API with rate limiting respect.
func makeRequest(endpoint string) (any, error) {
	url := fmt.Sprintf("%s%s", baseURL, endpoint)

	// Respect Jikan's rate limiting - 3 requests per second max
	time.Sleep(350 * time.Millisecond)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			slog.Error("failed to close response body", "error", closeErr)
		}
	}()

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("rate limited by Jikan API")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var result any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// GetAnimeByID retrieves anime information by MAL ID with 8-hour caching.
func GetAnimeByID(ctx context.Context, animeID int) (any, error) {
	return cache.NewCache("jikan_anime_by_id").TTL(8*time.Hour).Wrap(func() (any, error) {
		return makeRequest(fmt.Sprintf("/anime/%d", animeID))
	})(ctx, animeID)
}

// GetTopAnime gets top-rated anime with 2-hour caching (replaces "hot" anime).
func GetTopAnime(ctx context.Context) (any, error) {
	return cache.NewCache("jikan_top_anime").TTL(2 * time.Hour).Wrap(func() (any, error) {
		return makeRequest("/top/anime?type=tv&limit=20")
	})(ctx)
}

// GetSeasonalAnime gets current season anime with 4-hour caching.
func GetSeasonalAnime(ctx context.Context) (any, error) {
	return cache.NewCache("jikan_seasonal_anime").TTL(4 * time.Hour).Wrap(func() (any, error) {
		return makeRequest("/seasons/now")
	})(ctx)
}

// SearchAnime searches for anime by query with 4-hour caching.
func SearchAnime(ctx context.Context, query string, page int) (any, error) {
	return cache.NewCache("jikan_anime_search").TTL(4*time.Hour).Wrap(func() (any, error) {
		endpoint := fmt.Sprintf("/anime?q=%s", query)
		if page > 1 {
			endpoint += fmt.Sprintf("&page=%d", page)
		}
		return makeRequest(endpoint)
	})(ctx, query, page)
}

// GetAnimeRecommendations gets anime recommendations with 6-hour caching.
func GetAnimeRecommendations(ctx context.Context, animeID int) (any, error) {
	return cache.NewCache("jikan_anime_recommendations").TTL(6*time.Hour).Wrap(func() (any, error) {
		return makeRequest(fmt.Sprintf("/anime/%d/recommendations", animeID))
	})(ctx, animeID)
}

// GetRandomAnime gets a random anime recommendation with 2-hour caching.
func GetRandomAnime(ctx context.Context) (any, error) {
	return cache.NewCache("jikan_random_anime").TTL(2 * time.Hour).Wrap(func() (any, error) {
		return makeRequest("/random/anime")
	})(ctx)
}
