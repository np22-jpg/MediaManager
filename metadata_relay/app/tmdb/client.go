package tmdb

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"relay/app/cache"
)

var (
	baseURL string
	apiKey  string
)

// InitTMDB initializes the TMDB client with the provided API key and base URL.
// Logs a warning if the API key is not configured.
func InitTMDB(key, url string) {
	apiKey = key
	baseURL = url
	if apiKey == "" {
		fmt.Printf("WARNING: TMDB_API_KEY environment variable is not set. /tmdb/* endpoints will be dysfunctional.\n")
	}
}

// TrendingResponse represents the trending API response structure from TMDB.
type TrendingResponse struct {
	Page         int              `json:"page"`
	Results      []map[string]any `json:"results"`
	TotalPages   int              `json:"total_pages"`
	TotalResults int              `json:"total_results"`
}

// SearchResponse represents the search API response structure from TMDB.
type SearchResponse struct {
	Page         int              `json:"page"`
	Results      []map[string]any `json:"results"`
	TotalPages   int              `json:"total_pages"`
	TotalResults int              `json:"total_results"`
}

// makeRequest makes an HTTP request to TMDB API with proper error handling
// and response validation.
func makeRequest(endpoint string, params url.Values) (any, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("TMDB API key not configured")
	}

	params.Set("api_key", apiKey)
	url := fmt.Sprintf("%s%s?%s", baseURL, endpoint, params.Encode())

	resp, err := http.Get(url)
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

// GetTrendingTV gets trending TV shows from TMDB with caching for 2 hours.
func GetTrendingTV(ctx context.Context) (any, error) {
	return cache.NewCache("tmdb_tv_trending").TTL(2 * time.Hour).Wrap(func() (any, error) {
		params := url.Values{}
		return makeRequest("/trending/tv/week", params)
	})(ctx)
}

// SearchTV searches for TV shows with query and pagination, cached for 4 hours.
func SearchTV(ctx context.Context, query string, page int) (any, error) {
	return cache.NewCache("tmdb_tv_search").TTL(4*time.Hour).Wrap(func() (any, error) {
		params := url.Values{}
		params.Set("query", query)
		params.Set("page", strconv.Itoa(page))
		params.Set("include_adult", "true")
		return makeRequest("/search/tv", params)
	})(ctx, query, page)
}

// GetTVShow gets a specific TV show by ID with 12-hour caching.
func GetTVShow(ctx context.Context, showID int) (any, error) {
	return cache.NewCache("tmdb_tv_show").TTL(12*time.Hour).Wrap(func() (any, error) {
		params := url.Values{}
		return makeRequest(fmt.Sprintf("/tv/%d", showID), params)
	})(ctx, showID)
}

// GetTVSeason gets a specific season of a TV show with 8-hour caching.
func GetTVSeason(ctx context.Context, showID, seasonNumber int) (any, error) {
	return cache.NewCache("tmdb_tv_season").TTL(8*time.Hour).Wrap(func() (any, error) {
		params := url.Values{}
		return makeRequest(fmt.Sprintf("/tv/%d/season/%d", showID, seasonNumber), params)
	})(ctx, showID, seasonNumber)
}

// GetTrendingMovies gets trending movies from TMDB with 2-hour caching.
func GetTrendingMovies(ctx context.Context) (any, error) {
	return cache.NewCache("tmdb_movies_trending").TTL(2 * time.Hour).Wrap(func() (any, error) {
		params := url.Values{}
		return makeRequest("/trending/movie/week", params)
	})(ctx)
}

// SearchMovies searches for movies with query and pagination, cached for 4 hours.
func SearchMovies(ctx context.Context, query string, page int) (any, error) {
	return cache.NewCache("tmdb_movies_search").TTL(4*time.Hour).Wrap(func() (any, error) {
		params := url.Values{}
		params.Set("query", query)
		params.Set("page", strconv.Itoa(page))
		params.Set("include_adult", "true")
		return makeRequest("/search/movie", params)
	})(ctx, query, page)
}

// GetMovie gets a specific movie by ID with 8-hour caching.
func GetMovie(ctx context.Context, movieID int) (any, error) {
	return cache.NewCache("tmdb_movie").TTL(8*time.Hour).Wrap(func() (any, error) {
		params := url.Values{}
		return makeRequest(fmt.Sprintf("/movie/%d", movieID), params)
	})(ctx, movieID)
}
