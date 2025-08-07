package tvdb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"relay/app/cache"
	"relay/app/metrics"

	"github.com/caarlos0/env/v11"
)

const (
	baseURL  = "https://api4.thetvdb.com/v4"
	loginURL = "https://api4.thetvdb.com/v4/login"
)

type Config struct {
	APIKey string `env:"TVDB_API_KEY"`
}

var (
	apiKey      string
	accessToken string
	tokenExpiry time.Time
)

func init() {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		slog.Error("failed to parse TVDB configuration", "error", err)
	}

	apiKey = cfg.APIKey
	if apiKey == "" {
		fmt.Printf("WARNING: TVDB_API_KEY environment variable is not set\n")
	}
}

// LoginRequest represents the login request structure
type LoginRequest struct {
	APIKey string `json:"apikey"`
}

// LoginResponse represents the login response structure
type LoginResponse struct {
	Status string `json:"status"`
	Data   struct {
		Token string `json:"token"`
	} `json:"data"`
}

// authenticate gets a new access token from TVDB
func authenticate() error {
	if apiKey == "" {
		metrics.RecordAuthAttempt("tvdb", "failed")
		return fmt.Errorf("TVDB API key not configured")
	}

	loginReq := LoginRequest{APIKey: apiKey}
	jsonData, err := json.Marshal(loginReq)
	if err != nil {
		metrics.RecordAuthAttempt("tvdb", "failed")
		return fmt.Errorf("failed to marshal login request: %w", err)
	}

	resp, err := http.Post(loginURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		metrics.RecordAuthAttempt("tvdb", "failed")
		return fmt.Errorf("failed to authenticate: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			slog.Error("failed to close response body", "error", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		metrics.RecordAuthAttempt("tvdb", "failed")
		return fmt.Errorf("authentication failed with status %d", resp.StatusCode)
	}

	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		metrics.RecordAuthAttempt("tvdb", "failed")
		return fmt.Errorf("failed to decode login response: %w", err)
	}

	accessToken = loginResp.Data.Token
	tokenExpiry = time.Now().Add(23 * time.Hour) // TVDB tokens expire in 24h, refresh a bit early

	// Record successful authentication and update token expiry metrics
	metrics.RecordAuthAttempt("tvdb", "success")
	metrics.UpdateAuthTokenExpiry("tvdb", tokenExpiry)

	slog.Info("TVDB authentication successful",
		"expires_at", tokenExpiry.Format(time.RFC3339),
		"valid_for", time.Until(tokenExpiry).Round(time.Minute),
	)

	return nil
}

// ensureAuthenticated checks if we have a valid token and refreshes if needed
func ensureAuthenticated() error {
	if accessToken == "" || time.Now().After(tokenExpiry) {
		return authenticate()
	}
	return nil
}

// makeAuthenticatedRequest makes an HTTP request to TVDB API with authentication
func makeAuthenticatedRequest(endpoint string) (interface{}, error) {
	if err := ensureAuthenticated(); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s%s", baseURL, endpoint)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
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

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// GetTrendingTV gets all series (approximating trending functionality)
func GetTrendingTV(ctx context.Context) (interface{}, error) {
	return cache.NewCache("tvdb_tv_trending").TTL(2 * time.Hour).Wrap(func() (interface{}, error) {
		return makeAuthenticatedRequest("/series")
	})(ctx)
}

// SearchTV searches for TV shows
func SearchTV(ctx context.Context, query string) (interface{}, error) {
	return cache.NewCache("tvdb_tv_search").TTL(4*time.Hour).Wrap(func() (interface{}, error) {
		return makeAuthenticatedRequest(fmt.Sprintf("/search?query=%s", query))
	})(ctx, query)
}

// GetTVShow gets a specific TV show by ID
func GetTVShow(ctx context.Context, showID int) (interface{}, error) {
	return cache.NewCache("tvdb_tv_show").TTL(8*time.Hour).Wrap(func() (interface{}, error) {
		return makeAuthenticatedRequest(fmt.Sprintf("/series/%d/extended", showID))
	})(ctx, showID)
}

// GetTVSeason gets a specific season by ID
func GetTVSeason(ctx context.Context, seasonID int) (interface{}, error) {
	return cache.NewCache("tvdb_tv_season").TTL(8*time.Hour).Wrap(func() (interface{}, error) {
		return makeAuthenticatedRequest(fmt.Sprintf("/seasons/%d/extended", seasonID))
	})(ctx, seasonID)
}

// GetTrendingMovies gets all movies (approximating trending functionality)
func GetTrendingMovies(ctx context.Context) (interface{}, error) {
	return cache.NewCache("tvdb_movies_trending").TTL(2 * time.Hour).Wrap(func() (interface{}, error) {
		return makeAuthenticatedRequest("/movies")
	})(ctx)
}

// SearchMovies searches for movies
func SearchMovies(ctx context.Context, query string) (interface{}, error) {
	return cache.NewCache("tvdb_movies_search").TTL(4*time.Hour).Wrap(func() (interface{}, error) {
		return makeAuthenticatedRequest(fmt.Sprintf("/search?query=%s&type=movie", query))
	})(ctx, query)
}

// GetMovie gets a specific movie by ID
func GetMovie(ctx context.Context, movieID int) (interface{}, error) {
	return cache.NewCache("tvdb_movie").TTL(8*time.Hour).Wrap(func() (interface{}, error) {
		return makeAuthenticatedRequest(fmt.Sprintf("/movies/%d/extended", movieID))
	})(ctx, movieID)
}
