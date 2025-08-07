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

const (
	baseURL = "https://api.themoviedb.org/3"
)

var apiKey string

// initializes the TMDB client with the provided API key
func InitTMDB(key string) {
	apiKey = key
	if apiKey == "" {
		fmt.Printf("WARNING: TMDB_API_KEY environment variable is not set\n")
	}
}

// represents the trending API response structure
type TrendingResponse struct {
	Page         int              `json:"page"`
	Results      []map[string]any `json:"results"`
	TotalPages   int              `json:"total_pages"`
	TotalResults int              `json:"total_results"`
}

// represents the search API response structure
type SearchResponse struct {
	Page         int              `json:"page"`
	Results      []map[string]any `json:"results"`
	TotalPages   int              `json:"total_pages"`
	TotalResults int              `json:"total_results"`
}

// makeRequest makes an HTTP request to TMDB API
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

// gets trending TV shows
func GetTrendingTV(ctx context.Context) (any, error) {
	return cache.NewCache("tmdb_tv_trending").TTL(2 * time.Hour).Wrap(func() (any, error) {
		params := url.Values{}
		return makeRequest("/trending/tv/week", params)
	})(ctx)
}

// searches for TV shows
func SearchTV(ctx context.Context, query string, page int) (any, error) {
	return cache.NewCache("tmdb_tv_search").TTL(4*time.Hour).Wrap(func() (any, error) {
		params := url.Values{}
		params.Set("query", query)
		params.Set("page", strconv.Itoa(page))
		params.Set("include_adult", "true")
		return makeRequest("/search/tv", params)
	})(ctx, query, page)
}

// gets a specific TV show by ID
func GetTVShow(ctx context.Context, showID int) (any, error) {
	return cache.NewCache("tmdb_tv_show").TTL(12*time.Hour).Wrap(func() (any, error) {
		params := url.Values{}
		return makeRequest(fmt.Sprintf("/tv/%d", showID), params)
	})(ctx, showID)
}

// gets a specific season of a TV show
func GetTVSeason(ctx context.Context, showID, seasonNumber int) (any, error) {
	return cache.NewCache("tmdb_tv_season").TTL(8*time.Hour).Wrap(func() (any, error) {
		params := url.Values{}
		return makeRequest(fmt.Sprintf("/tv/%d/season/%d", showID, seasonNumber), params)
	})(ctx, showID, seasonNumber)
}

// gets trending movies
func GetTrendingMovies(ctx context.Context) (any, error) {
	return cache.NewCache("tmdb_movies_trending").TTL(2 * time.Hour).Wrap(func() (any, error) {
		params := url.Values{}
		return makeRequest("/trending/movie/week", params)
	})(ctx)
}

// searches for movies
func SearchMovies(ctx context.Context, query string, page int) (any, error) {
	return cache.NewCache("tmdb_movies_search").TTL(4*time.Hour).Wrap(func() (any, error) {
		params := url.Values{}
		params.Set("query", query)
		params.Set("page", strconv.Itoa(page))
		params.Set("include_adult", "true")
		return makeRequest("/search/movie", params)
	})(ctx, query, page)
}

// gets a specific movie by ID
func GetMovie(ctx context.Context, movieID int) (any, error) {
	return cache.NewCache("tmdb_movie").TTL(8*time.Hour).Wrap(func() (any, error) {
		params := url.Values{}
		return makeRequest(fmt.Sprintf("/movie/%d", movieID), params)
	})(ctx, movieID)
}
