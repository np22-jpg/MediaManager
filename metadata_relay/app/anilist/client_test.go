package anilist

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestAniListWithMockServer(t *testing.T) {
	// Create a test server that will handle all API requests
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Verify request headers
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected Accept: application/json, got %s", r.Header.Get("Accept"))
		}
		if r.Header.Get("User-Agent") == "" {
			t.Error("Expected User-Agent header to be set")
		}

		// Parse the GraphQL request
		var req GraphQLRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request: %v", err)
			w.WriteHeader(500)
			return
		}

		switch {
		case strings.Contains(req.Query, "Media(id: $id)"):
			// GetMediaByID request
			response := map[string]any{
				"data": map[string]any{
					"Media": map[string]any{
						"id": 123,
						"title": map[string]any{
							"romaji":  "Test Anime",
							"english": "Test Anime",
							"native":  "テストアニメ",
						},
						"type":        "ANIME",
						"format":      "TV",
						"status":      "FINISHED",
						"description": "A test anime description",
						"startDate": map[string]any{
							"year":  2024,
							"month": 1,
							"day":   1,
						},
						"episodes":     12,
						"duration":     24,
						"genres":       []string{"Action", "Adventure"},
						"averageScore": 85,
						"meanScore":    82,
						"popularity":   1000,
						"favourites":   500,
						"isAdult":      false,
						"coverImage": map[string]any{
							"extraLarge": "https://example.com/cover_xl.jpg",
							"large":      "https://example.com/cover_l.jpg",
							"medium":     "https://example.com/cover_m.jpg",
							"color":      "#ff6b35",
						},
						"bannerImage": "https://example.com/banner.jpg",
						"studios": map[string]any{
							"nodes": []map[string]any{
								{
									"id":                1,
									"name":              "Test Studio",
									"isAnimationStudio": true,
								},
							},
						},
						"tags": []map[string]any{
							{
								"id":          1,
								"name":        "Action",
								"description": "Action scenes",
								"category":    "Theme",
								"rank":        90,
								"isAdult":     false,
							},
						},
						"trailer": map[string]any{
							"id":        "abc123",
							"site":      "youtube",
							"thumbnail": "https://example.com/thumbnail.jpg",
						},
						"externalLinks": []map[string]any{
							{
								"id":       1,
								"url":      "https://example.com",
								"site":     "Official Site",
								"type":     "INFO",
								"language": "English",
								"color":    "#ffffff",
								"icon":     "https://example.com/icon.png",
							},
						},
						"nextAiringEpisode": map[string]any{
							"id":              1,
							"airingAt":        1642723200,
							"timeUntilAiring": 3600,
							"episode":         13,
							"mediaId":         123,
						},
					},
				},
			}
			_ = json.NewEncoder(w).Encode(response)

		case strings.Contains(req.Query, "media(search: $search"):
			// SearchMedia request
			response := map[string]any{
				"data": map[string]any{
					"Page": map[string]any{
						"pageInfo": map[string]any{
							"total":       1,
							"currentPage": 1,
							"lastPage":    1,
							"hasNextPage": false,
							"perPage":     20,
						},
						"media": []map[string]any{
							{
								"id": 456,
								"title": map[string]any{
									"romaji":  "Search Result",
									"english": "Search Result",
									"native":  "検索結果",
								},
								"type":        "ANIME",
								"format":      "TV",
								"status":      "RELEASING",
								"description": "A search result anime",
								"startDate": map[string]any{
									"year":  2024,
									"month": 4,
									"day":   1,
								},
								"episodes":     24,
								"duration":     24,
								"genres":       []string{"Drama", "Romance"},
								"averageScore": 75,
								"meanScore":    73,
								"popularity":   800,
								"favourites":   300,
								"isAdult":      false,
								"coverImage": map[string]any{
									"extraLarge": "https://example.com/search_cover_xl.jpg",
									"large":      "https://example.com/search_cover_l.jpg",
									"medium":     "https://example.com/search_cover_m.jpg",
									"color":      "#35a7ff",
								},
								"bannerImage": "https://example.com/search_banner.jpg",
								"studios": map[string]any{
									"nodes": []map[string]any{
										{
											"id":                2,
											"name":              "Search Studio",
											"isAnimationStudio": true,
										},
									},
								},
							},
						},
					},
				},
			}
			_ = json.NewEncoder(w).Encode(response)

		case strings.Contains(req.Query, "media(type: ANIME, sort: [TRENDING_DESC"):
			// GetTrendingAnime request
			response := map[string]any{
				"data": map[string]any{
					"Page": map[string]any{
						"pageInfo": map[string]any{
							"total":       2,
							"currentPage": 1,
							"lastPage":    1,
							"hasNextPage": false,
							"perPage":     20,
						},
						"media": []map[string]any{
							{
								"id": 789,
								"title": map[string]any{
									"romaji":  "Trending Anime 1",
									"english": "Trending Anime 1",
									"native":  "トレンドアニメ1",
								},
								"type":        "ANIME",
								"format":      "TV",
								"status":      "RELEASING",
								"description": "First trending anime",
								"startDate": map[string]any{
									"year": 2024,
								},
								"season":       "SPRING",
								"seasonYear":   2024,
								"episodes":     12,
								"duration":     24,
								"genres":       []string{"Action", "Fantasy"},
								"averageScore": 90,
								"meanScore":    88,
								"popularity":   5000,
								"favourites":   2500,
								"isAdult":      false,
								"coverImage": map[string]any{
									"extraLarge": "https://example.com/trending1_xl.jpg",
									"large":      "https://example.com/trending1_l.jpg",
									"medium":     "https://example.com/trending1_m.jpg",
									"color":      "#ff3535",
								},
								"bannerImage": "https://example.com/trending1_banner.jpg",
								"studios": map[string]any{
									"nodes": []map[string]any{
										{
											"id":                3,
											"name":              "Trending Studio",
											"isAnimationStudio": true,
										},
									},
								},
								"nextAiringEpisode": map[string]any{
									"id":              2,
									"airingAt":        1642809600,
									"timeUntilAiring": 86400,
									"episode":         8,
									"mediaId":         789,
								},
							},
							{
								"id": 790,
								"title": map[string]any{
									"romaji":  "Trending Anime 2",
									"english": "Trending Anime 2",
									"native":  "トレンドアニメ2",
								},
								"type":        "ANIME",
								"format":      "TV",
								"status":      "RELEASING",
								"description": "Second trending anime",
								"startDate": map[string]any{
									"year": 2024,
								},
								"season":       "SPRING",
								"seasonYear":   2024,
								"episodes":     24,
								"duration":     24,
								"genres":       []string{"Comedy", "Slice of Life"},
								"averageScore": 85,
								"meanScore":    83,
								"popularity":   4500,
								"favourites":   2200,
								"isAdult":      false,
								"coverImage": map[string]any{
									"extraLarge": "https://example.com/trending2_xl.jpg",
									"large":      "https://example.com/trending2_l.jpg",
									"medium":     "https://example.com/trending2_m.jpg",
									"color":      "#35ff35",
								},
								"bannerImage": "https://example.com/trending2_banner.jpg",
								"studios": map[string]any{
									"nodes": []map[string]any{
										{
											"id":                4,
											"name":              "Comedy Studio",
											"isAnimationStudio": true,
										},
									},
								},
								"nextAiringEpisode": map[string]any{
									"id":              3,
									"airingAt":        1642723200,
									"timeUntilAiring": 3600,
									"episode":         5,
									"mediaId":         790,
								},
							},
						},
					},
				},
			}
			_ = json.NewEncoder(w).Encode(response)

		case strings.Contains(req.Query, "seasonYear: $year, season: $season"):
			// GetSeasonalAnime request
			response := map[string]any{
				"data": map[string]any{
					"Page": map[string]any{
						"pageInfo": map[string]any{
							"total":       1,
							"currentPage": 1,
							"lastPage":    1,
							"hasNextPage": false,
							"perPage":     20,
						},
						"media": []map[string]any{
							{
								"id": 999,
								"title": map[string]any{
									"romaji":  "Seasonal Anime",
									"english": "Seasonal Anime",
									"native":  "シーズナルアニメ",
								},
								"type":        "ANIME",
								"format":      "TV",
								"status":      "FINISHED",
								"description": "A seasonal anime",
								"startDate": map[string]any{
									"year":  2024,
									"month": 1,
									"day":   1,
								},
								"season":       "WINTER",
								"seasonYear":   2024,
								"episodes":     12,
								"duration":     24,
								"genres":       []string{"Drama"},
								"averageScore": 80,
								"meanScore":    78,
								"popularity":   1500,
								"favourites":   750,
								"isAdult":      false,
								"coverImage": map[string]any{
									"extraLarge": "https://example.com/seasonal_xl.jpg",
									"large":      "https://example.com/seasonal_l.jpg",
									"medium":     "https://example.com/seasonal_m.jpg",
									"color":      "#3535ff",
								},
								"bannerImage": "https://example.com/seasonal_banner.jpg",
								"studios": map[string]any{
									"nodes": []map[string]any{
										{
											"id":                5,
											"name":              "Seasonal Studio",
											"isAnimationStudio": true,
										},
									},
								},
								"nextAiringEpisode": nil,
							},
						},
					},
				},
			}
			_ = json.NewEncoder(w).Encode(response)

		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	// Store original env vars and set test values
	origURL := os.Getenv("ANILIST_GRAPHQL_URL")
	origUA := os.Getenv("ANILIST_USER_AGENT")
	if err := os.Setenv("ANILIST_GRAPHQL_URL", server.URL); err != nil {
		t.Fatalf("Failed to set ANILIST_GRAPHQL_URL: %v", err)
	}
	if err := os.Setenv("ANILIST_USER_AGENT", "Test-Agent/1.0"); err != nil {
		t.Fatalf("Failed to set ANILIST_USER_AGENT: %v", err)
	}
	defer func() {
		_ = os.Setenv("ANILIST_GRAPHQL_URL", origURL)
		_ = os.Setenv("ANILIST_USER_AGENT", origUA)
	}()

	ctx := context.Background()

	t.Run("GetMediaByID", func(t *testing.T) {
		media, err := GetMediaByID(ctx, 123)
		if err != nil {
			t.Fatalf("GetMediaByID failed: %v", err)
		}
		if media == nil {
			t.Fatal("Media should not be nil")
		}

		if media.ID != 123 {
			t.Errorf("Expected media ID 123, got %d", media.ID)
		}
		if media.Title.Romaji != "Test Anime" {
			t.Errorf("Expected title 'Test Anime', got %s", media.Title.Romaji)
		}
		if media.Type != "ANIME" {
			t.Errorf("Expected type 'ANIME', got %s", media.Type)
		}
		if media.Episodes != 12 {
			t.Errorf("Expected 12 episodes, got %d", media.Episodes)
		}
		if len(media.Genres) != 2 {
			t.Errorf("Expected 2 genres, got %d", len(media.Genres))
		}
		if media.AverageScore != 85 {
			t.Errorf("Expected average score 85, got %d", media.AverageScore)
		}
	})

	t.Run("SearchMedia", func(t *testing.T) {
		page, err := SearchMedia(ctx, "test", "ANIME", 1, 20)
		if err != nil {
			t.Fatalf("SearchMedia failed: %v", err)
		}
		if page == nil {
			t.Fatal("Page should not be nil")
		}

		if page.PageInfo.Total != 1 {
			t.Errorf("Expected total 1, got %d", page.PageInfo.Total)
		}
		if len(page.Media) != 1 {
			t.Errorf("Expected 1 media item, got %d", len(page.Media))
		}

		media := page.Media[0]
		if media.ID != 456 {
			t.Errorf("Expected media ID 456, got %d", media.ID)
		}
		if media.Title.Romaji != "Search Result" {
			t.Errorf("Expected title 'Search Result', got %s", media.Title.Romaji)
		}
	})

	t.Run("GetTrendingAnime", func(t *testing.T) {
		page, err := GetTrendingAnime(ctx, 1, 20)
		if err != nil {
			t.Fatalf("GetTrendingAnime failed: %v", err)
		}
		if page == nil {
			t.Fatal("Page should not be nil")
		}

		if page.PageInfo.Total != 2 {
			t.Errorf("Expected total 2, got %d", page.PageInfo.Total)
		}
		if len(page.Media) != 2 {
			t.Errorf("Expected 2 media items, got %d", len(page.Media))
		}

		// Check first trending anime
		media1 := page.Media[0]
		if media1.ID != 789 {
			t.Errorf("Expected media ID 789, got %d", media1.ID)
		}
		if media1.Title.Romaji != "Trending Anime 1" {
			t.Errorf("Expected title 'Trending Anime 1', got %s", media1.Title.Romaji)
		}
		if media1.Season != "SPRING" {
			t.Errorf("Expected season 'SPRING', got %s", media1.Season)
		}
		if media1.SeasonYear != 2024 {
			t.Errorf("Expected season year 2024, got %d", media1.SeasonYear)
		}

		// Check second trending anime
		media2 := page.Media[1]
		if media2.ID != 790 {
			t.Errorf("Expected media ID 790, got %d", media2.ID)
		}
		if media2.Title.Romaji != "Trending Anime 2" {
			t.Errorf("Expected title 'Trending Anime 2', got %s", media2.Title.Romaji)
		}
	})

	t.Run("GetSeasonalAnime", func(t *testing.T) {
		page, err := GetSeasonalAnime(ctx, 2024, "WINTER", 1, 20)
		if err != nil {
			t.Fatalf("GetSeasonalAnime failed: %v", err)
		}
		if page == nil {
			t.Fatal("Page should not be nil")
		}

		if page.PageInfo.Total != 1 {
			t.Errorf("Expected total 1, got %d", page.PageInfo.Total)
		}
		if len(page.Media) != 1 {
			t.Errorf("Expected 1 media item, got %d", len(page.Media))
		}

		media := page.Media[0]
		if media.ID != 999 {
			t.Errorf("Expected media ID 999, got %d", media.ID)
		}
		if media.Title.Romaji != "Seasonal Anime" {
			t.Errorf("Expected title 'Seasonal Anime', got %s", media.Title.Romaji)
		}
		if media.Season != "WINTER" {
			t.Errorf("Expected season 'WINTER', got %s", media.Season)
		}
		if media.SeasonYear != 2024 {
			t.Errorf("Expected season year 2024, got %d", media.SeasonYear)
		}
	})

	t.Run("DefaultParameters", func(t *testing.T) {
		// Test with default parameters (0 values should be converted to defaults)
		page, err := SearchMedia(ctx, "test", "", 0, 0)
		if err != nil {
			t.Fatalf("SearchMedia with defaults failed: %v", err)
		}
		if page == nil {
			t.Fatal("Page should not be nil")
		}
		// Should still return results with default page=1, perPage=20
	})
}

func TestAniListErrorHandling(t *testing.T) {
	// Test with server that returns GraphQL errors
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := map[string]any{
			"data": nil,
			"errors": []map[string]any{
				{
					"message": "Test GraphQL error",
					"path":    []string{"Media"},
				},
			},
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Store original env var and set test value
	origURL := os.Getenv("ANILIST_GRAPHQL_URL")
	_ = os.Setenv("ANILIST_GRAPHQL_URL", server.URL)
	defer func() {
		_ = os.Setenv("ANILIST_GRAPHQL_URL", origURL)
	}()

	ctx := context.Background()

	t.Run("GraphQLError", func(t *testing.T) {
		media, err := GetMediaByID(ctx, 123)
		if err == nil {
			t.Error("Expected error but got none")
		}
		if media != nil {
			t.Error("Expected nil media but got non-nil")
		}
		if !strings.Contains(err.Error(), "GraphQL errors") {
			t.Errorf("Expected error to contain 'GraphQL errors', got: %v", err)
		}
	})
}

func TestAniListServerErrors(t *testing.T) {
	// Test with server that returns HTTP errors
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		_, _ = w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	// Store original env var and set test value
	origURL := os.Getenv("ANILIST_GRAPHQL_URL")
	_ = os.Setenv("ANILIST_GRAPHQL_URL", server.URL)
	defer func() {
		_ = os.Setenv("ANILIST_GRAPHQL_URL", origURL)
	}()

	ctx := context.Background()

	t.Run("ServerError", func(t *testing.T) {
		media, err := GetMediaByID(ctx, 123)
		if err == nil {
			t.Error("Expected error but got none")
		}
		if media != nil {
			t.Error("Expected nil media but got non-nil")
		}
		if !strings.Contains(err.Error(), "API request failed with status 500") {
			t.Errorf("Expected error to contain status 500, got: %v", err)
		}
	})
}

func TestAniListConfig(t *testing.T) {
	t.Run("DefaultConfig", func(t *testing.T) {
		// Clear env vars to test defaults
		origURL := os.Getenv("ANILIST_GRAPHQL_URL")
		origUA := os.Getenv("ANILIST_USER_AGENT")
		if err := os.Unsetenv("ANILIST_GRAPHQL_URL"); err != nil {
			t.Fatalf("Failed to unset ANILIST_GRAPHQL_URL: %v", err)
		}
		if err := os.Unsetenv("ANILIST_USER_AGENT"); err != nil {
			t.Fatalf("Failed to unset ANILIST_USER_AGENT: %v", err)
		}
		defer func() {
			_ = os.Setenv("ANILIST_GRAPHQL_URL", origURL)
			_ = os.Setenv("ANILIST_USER_AGENT", origUA)
		}()

		url, ua := getConfig()
		if url != "https://graphql.anilist.co" {
			t.Errorf("Expected default URL 'https://graphql.anilist.co', got %s", url)
		}
		if ua != "MediaManager-Relay/1.0" {
			t.Errorf("Expected default User-Agent 'MediaManager-Relay/1.0', got %s", ua)
		}
	})

	t.Run("CustomConfig", func(t *testing.T) {
		origURL := os.Getenv("ANILIST_GRAPHQL_URL")
		origUA := os.Getenv("ANILIST_USER_AGENT")
		_ = os.Setenv("ANILIST_GRAPHQL_URL", "https://custom.anilist.api")
		_ = os.Setenv("ANILIST_USER_AGENT", "Custom-Agent/2.0")
		defer func() {
			_ = os.Setenv("ANILIST_GRAPHQL_URL", origURL)
			_ = os.Setenv("ANILIST_USER_AGENT", origUA)
		}()

		url, ua := getConfig()
		if url != "https://custom.anilist.api" {
			t.Errorf("Expected custom URL 'https://custom.anilist.api', got %s", url)
		}
		if ua != "Custom-Agent/2.0" {
			t.Errorf("Expected custom User-Agent 'Custom-Agent/2.0', got %s", ua)
		}
	})
}
