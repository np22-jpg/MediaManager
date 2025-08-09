package anilist

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"relay/app/cache"
)

// getConfig returns the AniList configuration values from environment variables
func getConfig() (anilistGraphQLURL, userAgent string) {
	anilistGraphQLURL = os.Getenv("ANILIST_GRAPHQL_URL")
	if anilistGraphQLURL == "" {
		anilistGraphQLURL = "https://graphql.anilist.co"
	}

	userAgent = os.Getenv("ANILIST_USER_AGENT")
	if userAgent == "" {
		userAgent = "MediaManager-Relay/1.0"
	}

	return
}

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

// GraphQLRequest represents a GraphQL request
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// GraphQLResponse represents a GraphQL response
type GraphQLResponse struct {
	Data   interface{} `json:"data"`
	Errors []struct {
		Message string        `json:"message"`
		Path    []interface{} `json:"path,omitempty"`
	} `json:"errors,omitempty"`
}

// Media represents an anime/manga from AniList
type Media struct {
	ID                int              `json:"id"`
	Title             MediaTitle       `json:"title"`
	Type              string           `json:"type"`
	Format            string           `json:"format"`
	Status            string           `json:"status"`
	Description       string           `json:"description"`
	StartDate         *FuzzyDate       `json:"startDate"`
	EndDate           *FuzzyDate       `json:"endDate"`
	Season            string           `json:"season"`
	SeasonYear        int              `json:"seasonYear"`
	Episodes          int              `json:"episodes"`
	Duration          int              `json:"duration"`
	Chapters          int              `json:"chapters"`
	Volumes           int              `json:"volumes"`
	Genres            []string         `json:"genres"`
	AverageScore      int              `json:"averageScore"`
	MeanScore         int              `json:"meanScore"`
	Popularity        int              `json:"popularity"`
	Favourites        int              `json:"favourites"`
	IsAdult           bool             `json:"isAdult"`
	CoverImage        CoverImage       `json:"coverImage"`
	BannerImage       string           `json:"bannerImage"`
	Studios           StudioConnection `json:"studios"`
	Tags              []MediaTag       `json:"tags"`
	Trailer           *MediaTrailer    `json:"trailer"`
	ExternalLinks     []ExternalLink   `json:"externalLinks"`
	NextAiringEpisode *AiringSchedule  `json:"nextAiringEpisode"`
}

type MediaTitle struct {
	Romaji  string `json:"romaji"`
	English string `json:"english"`
	Native  string `json:"native"`
}

type FuzzyDate struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

type CoverImage struct {
	ExtraLarge string `json:"extraLarge"`
	Large      string `json:"large"`
	Medium     string `json:"medium"`
	Color      string `json:"color"`
}

type MediaTrailer struct {
	ID        string `json:"id"`
	Site      string `json:"site"`
	Thumbnail string `json:"thumbnail"`
}

type MediaTag struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Rank        int    `json:"rank"`
	IsAdult     bool   `json:"isAdult"`
}

type StudioConnection struct {
	Nodes []Studio `json:"nodes"`
}

type Studio struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	IsAnimationStudio bool   `json:"isAnimationStudio"`
}

type ExternalLink struct {
	ID       int    `json:"id"`
	URL      string `json:"url"`
	Site     string `json:"site"`
	Type     string `json:"type"`
	Language string `json:"language"`
	Color    string `json:"color"`
	Icon     string `json:"icon"`
}

type AiringSchedule struct {
	ID              int `json:"id"`
	AiringAt        int `json:"airingAt"`
	TimeUntilAiring int `json:"timeUntilAiring"`
	Episode         int `json:"episode"`
	MediaID         int `json:"mediaId"`
}

// Page represents paginated results
type Page struct {
	PageInfo PageInfo `json:"pageInfo"`
	Media    []Media  `json:"media"`
}

type PageInfo struct {
	Total       int  `json:"total"`
	CurrentPage int  `json:"currentPage"`
	LastPage    int  `json:"lastPage"`
	HasNextPage bool `json:"hasNextPage"`
	PerPage     int  `json:"perPage"`
}

// makeGraphQLRequest executes a GraphQL request to AniList API
func makeGraphQLRequest(ctx context.Context, query string, variables map[string]interface{}) (*GraphQLResponse, error) {
	anilistGraphQLURL, userAgent := getConfig()

	reqBody := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", anilistGraphQLURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			slog.Error("failed to close response body", "error", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var gqlResp GraphQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&gqlResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(gqlResp.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL errors: %v", gqlResp.Errors)
	}

	return &gqlResp, nil
}

// GetMediaByID retrieves anime/manga by AniList ID with 8-hour caching
func GetMediaByID(ctx context.Context, mediaID int) (*Media, error) {
	result, err := cache.NewCache("anilist_media_by_id").TTL(8*time.Hour).Wrap(func() (interface{}, error) {
		query := `
		query ($id: Int) {
			Media(id: $id) {
				id
				title {
					romaji
					english
					native
				}
				type
				format
				status
				description
				startDate {
					year
					month
					day
				}
				endDate {
					year
					month
					day
				}
				season
				seasonYear
				episodes
				duration
				chapters
				volumes
				genres
				averageScore
				meanScore
				popularity
				favourites
				isAdult
				coverImage {
					extraLarge
					large
					medium
					color
				}
				bannerImage
				studios {
					nodes {
						id
						name
						isAnimationStudio
					}
				}
				tags {
					id
					name
					description
					category
					rank
					isAdult
				}
				trailer {
					id
					site
					thumbnail
				}
				externalLinks {
					id
					url
					site
					type
					language
					color
					icon
				}
				nextAiringEpisode {
					id
					airingAt
					timeUntilAiring
					episode
					mediaId
				}
			}
		}`

		variables := map[string]interface{}{
			"id": mediaID,
		}

		resp, err := makeGraphQLRequest(ctx, query, variables)
		if err != nil {
			return nil, err
		}

		// Extract Media from response
		dataMap, ok := resp.Data.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected response format")
		}

		mediaData, ok := dataMap["Media"]
		if !ok || mediaData == nil {
			return nil, fmt.Errorf("media not found")
		}

		// Convert to JSON and back to struct for proper type conversion
		jsonData, err := json.Marshal(mediaData)
		if err != nil {
			return nil, err
		}

		var media Media
		if err := json.Unmarshal(jsonData, &media); err != nil {
			return nil, err
		}

		return &media, nil
	})(ctx, mediaID)

	if err != nil {
		return nil, err
	}

	return result.(*Media), nil
}

// SearchMedia searches for anime/manga by query with 4-hour caching
func SearchMedia(ctx context.Context, search string, mediaType string, page int, perPage int) (*Page, error) {
	if perPage == 0 {
		perPage = 20
	}
	if page == 0 {
		page = 1
	}

	result, err := cache.NewCache("anilist_search").TTL(4*time.Hour).Wrap(func() (interface{}, error) {
		query := `
		query ($search: String, $type: MediaType, $page: Int, $perPage: Int) {
			Page(page: $page, perPage: $perPage) {
				pageInfo {
					total
					currentPage
					lastPage
					hasNextPage
					perPage
				}
				media(search: $search, type: $type, sort: [POPULARITY_DESC, SCORE_DESC]) {
					id
					title {
						romaji
						english
						native
					}
					type
					format
					status
					description
					startDate {
						year
						month
						day
					}
					season
					seasonYear
					episodes
					duration
					chapters
					volumes
					genres
					averageScore
					meanScore
					popularity
					favourites
					isAdult
					coverImage {
						extraLarge
						large
						medium
						color
					}
					bannerImage
					studios {
						nodes {
							id
							name
							isAnimationStudio
						}
					}
				}
			}
		}`

		variables := map[string]interface{}{
			"search":  search,
			"page":    page,
			"perPage": perPage,
		}

		if mediaType != "" {
			variables["type"] = mediaType
		}

		resp, err := makeGraphQLRequest(ctx, query, variables)
		if err != nil {
			return nil, err
		}

		// Extract Page from response
		dataMap, ok := resp.Data.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected response format")
		}

		pageData, ok := dataMap["Page"]
		if !ok {
			return nil, fmt.Errorf("page data not found")
		}

		// Convert to JSON and back to struct for proper type conversion
		jsonData, err := json.Marshal(pageData)
		if err != nil {
			return nil, err
		}

		var pageResult Page
		if err := json.Unmarshal(jsonData, &pageResult); err != nil {
			return nil, err
		}

		return &pageResult, nil
	})(ctx, search, mediaType, page, perPage)

	if err != nil {
		return nil, err
	}

	return result.(*Page), nil
}

// GetTrendingAnime gets trending anime with 2-hour caching
func GetTrendingAnime(ctx context.Context, page int, perPage int) (*Page, error) {
	if perPage == 0 {
		perPage = 20
	}
	if page == 0 {
		page = 1
	}

	result, err := cache.NewCache("anilist_trending_anime").TTL(2*time.Hour).Wrap(func() (interface{}, error) {
		query := `
		query ($page: Int, $perPage: Int) {
			Page(page: $page, perPage: $perPage) {
				pageInfo {
					total
					currentPage
					lastPage
					hasNextPage
					perPage
				}
				media(type: ANIME, sort: [TRENDING_DESC, POPULARITY_DESC]) {
					id
					title {
						romaji
						english
						native
					}
					type
					format
					status
					description
					startDate {
						year
					}
					season
					seasonYear
					episodes
					duration
					genres
					averageScore
					meanScore
					popularity
					favourites
					isAdult
					coverImage {
						extraLarge
						large
						medium
						color
					}
					bannerImage
					studios {
						nodes {
							id
							name
							isAnimationStudio
						}
					}
					nextAiringEpisode {
						id
						airingAt
						timeUntilAiring
						episode
						mediaId
					}
				}
			}
		}`

		variables := map[string]interface{}{
			"page":    page,
			"perPage": perPage,
		}

		resp, err := makeGraphQLRequest(ctx, query, variables)
		if err != nil {
			return nil, err
		}

		// Extract Page from response
		dataMap, ok := resp.Data.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected response format")
		}

		pageData, ok := dataMap["Page"]
		if !ok {
			return nil, fmt.Errorf("page data not found")
		}

		// Convert to JSON and back to struct
		jsonData, err := json.Marshal(pageData)
		if err != nil {
			return nil, err
		}

		var pageResult Page
		if err := json.Unmarshal(jsonData, &pageResult); err != nil {
			return nil, err
		}

		return &pageResult, nil
	})(ctx, page, perPage)

	if err != nil {
		return nil, err
	}

	return result.(*Page), nil
}

// GetSeasonalAnime gets current season anime with 4-hour caching
func GetSeasonalAnime(ctx context.Context, year int, season string, page int, perPage int) (*Page, error) {
	if perPage == 0 {
		perPage = 20
	}
	if page == 0 {
		page = 1
	}

	// If year/season not specified, use current
	if year == 0 {
		now := time.Now()
		year = now.Year()
		month := now.Month()
		switch {
		case month >= 12 || month <= 2:
			season = "WINTER"
		case month >= 3 && month <= 5:
			season = "SPRING"
		case month >= 6 && month <= 8:
			season = "SUMMER"
		case month >= 9 && month <= 11:
			season = "FALL"
		}
	}

	result, err := cache.NewCache("anilist_seasonal_anime").TTL(4*time.Hour).Wrap(func() (interface{}, error) {
		query := `
		query ($year: Int, $season: MediaSeason, $page: Int, $perPage: Int) {
			Page(page: $page, perPage: $perPage) {
				pageInfo {
					total
					currentPage
					lastPage
					hasNextPage
					perPage
				}
				media(type: ANIME, seasonYear: $year, season: $season, sort: [POPULARITY_DESC]) {
					id
					title {
						romaji
						english
						native
					}
					type
					format
					status
					description
					startDate {
						year
						month
						day
					}
					season
					seasonYear
					episodes
					duration
					genres
					averageScore
					meanScore
					popularity
					favourites
					isAdult
					coverImage {
						extraLarge
						large
						medium
						color
					}
					bannerImage
					studios {
						nodes {
							id
							name
							isAnimationStudio
						}
					}
					nextAiringEpisode {
						id
						airingAt
						timeUntilAiring
						episode
						mediaId
					}
				}
			}
		}`

		variables := map[string]interface{}{
			"year":    year,
			"season":  season,
			"page":    page,
			"perPage": perPage,
		}

		resp, err := makeGraphQLRequest(ctx, query, variables)
		if err != nil {
			return nil, err
		}

		// Extract Page from response
		dataMap, ok := resp.Data.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected response format")
		}

		pageData, ok := dataMap["Page"]
		if !ok {
			return nil, fmt.Errorf("page data not found")
		}

		// Convert to JSON and back to struct
		jsonData, err := json.Marshal(pageData)
		if err != nil {
			return nil, err
		}

		var pageResult Page
		if err := json.Unmarshal(jsonData, &pageResult); err != nil {
			return nil, err
		}

		return &pageResult, nil
	})(ctx, year, season, page, perPage)

	if err != nil {
		return nil, err
	}

	return result.(*Page), nil
}
