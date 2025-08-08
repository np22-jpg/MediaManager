package anidb

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"relay/app/cache"
)

var (
	baseURL     string
	clientName  string
	clientVer   string
	protocolVer string
)

// InitAniDB initializes the AniDB client with the provided configuration.
func InitAniDB(url, client, clientVersion string) {
	baseURL = url
	clientName = client
	clientVer = clientVersion
	protocolVer = "1"

	if baseURL == "" {
		baseURL = "http://api.anidb.net:9001/httpapi"
	}
	if clientName == "" {
		clientName = "metadatarelay"
	}
	if clientVer == "" {
		clientVer = "1"
	}

	slog.Info("AniDB initialized", "baseUrl", baseURL, "client", clientName, "version", clientVer)
}

// AnimeInfo represents anime information from AniDB.
type AnimeInfo struct {
	XMLName      xml.Name           `xml:"anime" json:"-"`
	ID           int                `xml:"id,attr" json:"id"`
	Restricted   bool               `xml:"restricted,attr" json:"restricted"`
	Type         string             `xml:"type" json:"type"`
	EpisodeCount int                `xml:"episodecount" json:"episode_count"`
	StartDate    string             `xml:"startdate" json:"start_date"`
	EndDate      string             `xml:"enddate" json:"end_date"`
	Titles       Titles             `xml:"titles" json:"titles"`
	RelatedAnime []RelatedAnime     `xml:"relatedanime>anime" json:"related_anime"`
	SimilarAnime []SimilarAnimePair `xml:"similaranime>anime" json:"similar_anime"`
	Description  string             `xml:"description" json:"description"`
	Ratings      Ratings            `xml:"ratings" json:"ratings"`
	Picture      string             `xml:"picture" json:"picture"`
	URL          string             `xml:"url" json:"url"`
	Tags         []Tag              `xml:"tags>tag" json:"tags"`
	Characters   []Character        `xml:"characters>character" json:"characters"`
	Episodes     []Episode          `xml:"episodes>episode" json:"episodes"`
}

// Titles represents anime titles.
type Titles struct {
	Title []Title `xml:"title" json:"titles"`
}

// Title represents a single anime title.
type Title struct {
	Lang  string `xml:"lang,attr" json:"lang"`
	Type  string `xml:"type,attr" json:"type"`
	Value string `xml:",chardata" json:"value"`
}

// RelatedAnime represents related anime information.
type RelatedAnime struct {
	ID    int    `xml:"id,attr" json:"id"`
	Type  string `xml:"type,attr" json:"type"`
	Title string `xml:",chardata" json:"title"`
}

// SimilarAnimeSource represents the source anime in a similar anime pair.
type SimilarAnimeSource struct {
	ID         int    `xml:"aid,attr" json:"id"`
	Restricted bool   `xml:"restricted,attr" json:"restricted"`
	Title      Title  `xml:"title" json:"title"`
	Picture    string `xml:"picture" json:"picture,omitempty"`
}

// SimilarAnimeTarget represents the target anime in a similar anime pair.
type SimilarAnimeTarget struct {
	ID         int    `xml:"aid,attr" json:"id"`
	Restricted bool   `xml:"restricted,attr" json:"restricted"`
	Title      Title  `xml:"title" json:"title"`
	Picture    string `xml:"picture" json:"picture,omitempty"`
}

// SimilarAnimePair represents a pair of similar anime.
type SimilarAnimePair struct {
	Source SimilarAnimeSource `xml:"source" json:"source"`
	Target SimilarAnimeTarget `xml:"target" json:"target"`
}

// Ratings represents anime ratings.
type Ratings struct {
	Permanent Rating `xml:"permanent" json:"permanent"`
	Temporary Rating `xml:"temporary" json:"temporary"`
	Review    Rating `xml:"review" json:"review"`
}

// Rating represents a single rating.
type Rating struct {
	Count int     `xml:"count,attr" json:"count"`
	Value float64 `xml:",chardata" json:"value"`
}

// Tag represents an anime tag.
type Tag struct {
	ID            int    `xml:"id,attr" json:"id"`
	ParentID      int    `xml:"parentid,attr" json:"parent_id"`
	Weight        int    `xml:"weight,attr" json:"weight"`
	LocalSpoiler  bool   `xml:"localspoiler,attr" json:"local_spoiler"`
	GlobalSpoiler bool   `xml:"globalspoiler,attr" json:"global_spoiler"`
	Verified      bool   `xml:"verified,attr" json:"verified"`
	Update        string `xml:"update,attr" json:"update"`
	Name          string `xml:"name" json:"name"`
	Description   string `xml:"description" json:"description"`
	PictureURL    string `xml:"picurl" json:"picture_url"`
}

// Character represents an anime character.
type Character struct {
	ID          int    `xml:"id,attr" json:"id"`
	Type        string `xml:"type,attr" json:"type"`
	Update      string `xml:"update,attr" json:"update"`
	Rating      Rating `xml:"rating" json:"rating"`
	Name        string `xml:"name" json:"name"`
	Gender      string `xml:"gender" json:"gender"`
	Description string `xml:"description" json:"description"`
	Picture     string `xml:"picture" json:"picture"`
}

// Episode represents an anime episode.
type Episode struct {
	ID      int    `xml:"id,attr" json:"id"`
	Update  string `xml:"update,attr" json:"update"`
	EpNo    EpNo   `xml:"epno" json:"epno"`
	Length  int    `xml:"length" json:"length"`
	AirDate string `xml:"airdate" json:"air_date"`
	Rating  Rating `xml:"rating" json:"rating"`
	Title   Title  `xml:"title" json:"title"`
}

// EpNo represents episode number information.
type EpNo struct {
	Type  int `xml:"type,attr" json:"type"`
	Value int `xml:",chardata" json:"value"`
}

// ErrorResponse represents an AniDB error response.
type ErrorResponse struct {
	XMLName xml.Name `xml:"error"`
	Message string   `xml:",chardata"`
}

// makeRequest makes an HTTP request to AniDB API with proper error handling.
func makeRequest(request string, params url.Values) (any, error) {
	if clientName == "" {
		return nil, fmt.Errorf("AniDB client not configured")
	}

	// Build base parameters
	reqParams := url.Values{}
	reqParams.Set("client", clientName)
	reqParams.Set("clientver", clientVer)
	reqParams.Set("protover", protocolVer)
	reqParams.Set("request", request)

	// Add additional parameters
	for k, v := range params {
		reqParams[k] = v
	}

	reqURL := fmt.Sprintf("%s?%s", baseURL, reqParams.Encode())

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

	// Try to decode as error first
	var errorResp ErrorResponse
	if err := xml.NewDecoder(resp.Body).Decode(&errorResp); err == nil && errorResp.Message != "" {
		return nil, fmt.Errorf("AniDB error: %s", errorResp.Message)
	}

	// Reset response body and decode as requested type
	resp, err = http.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to re-make request: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			slog.Error("failed to close response body", "error", closeErr)
		}
	}()

	var result any
	switch request {
	case "anime":
		var anime AnimeInfo
		if err := xml.NewDecoder(resp.Body).Decode(&anime); err != nil {
			return nil, fmt.Errorf("failed to decode anime response: %w", err)
		}
		result = anime
	case "hotanime":
		var hotAnime struct {
			XMLName xml.Name    `xml:"hotanime" json:"-"`
			Anime   []AnimeInfo `xml:"anime" json:"anime"`
		}
		if err := xml.NewDecoder(resp.Body).Decode(&hotAnime); err != nil {
			return nil, fmt.Errorf("failed to decode hot anime response: %w", err)
		}
		result = hotAnime.Anime
	case "randomrecommendation":
		var recommendation struct {
			XMLName        xml.Name `xml:"randomrecommendation" json:"-"`
			Recommendation struct {
				Anime AnimeInfo `xml:"anime" json:"anime"`
			} `xml:"recommendation" json:"recommendation"`
		}
		if err := xml.NewDecoder(resp.Body).Decode(&recommendation); err != nil {
			return nil, fmt.Errorf("failed to decode recommendation response: %w", err)
		}
		result = recommendation.Recommendation.Anime
	case "randomsimilar":
		var similar struct {
			XMLName xml.Name         `xml:"randomsimilar" json:"-"`
			Similar SimilarAnimePair `xml:"similar" json:"similar"`
		}
		if err := xml.NewDecoder(resp.Body).Decode(&similar); err != nil {
			return nil, fmt.Errorf("failed to decode similar response: %w", err)
		}
		result = similar.Similar
	case "main":
		var main struct {
			XMLName  xml.Name `xml:"main" json:"-"`
			HotAnime struct {
				Anime []AnimeInfo `xml:"anime" json:"anime"`
			} `xml:"hotanime" json:"hotanime"`
			Recommendation struct {
				Recommendation struct {
					Anime AnimeInfo `xml:"anime" json:"anime"`
				} `xml:"recommendation" json:"recommendation"`
			} `xml:"randomrecommendation" json:"randomrecommendation"`
		}
		if err := xml.NewDecoder(resp.Body).Decode(&main); err != nil {
			return nil, fmt.Errorf("failed to decode main response: %w", err)
		}
		result = map[string]any{
			"hotanime":    main.HotAnime.Anime,
			"recommended": main.Recommendation.Recommendation.Anime,
		}
	default:
		// For unknown request types, return raw XML as string
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}
		result = string(body)
	}

	return result, nil
}

// GetAnimeByID retrieves anime information by AniDB ID with 8-hour caching.
func GetAnimeByID(ctx context.Context, animeID int) (any, error) {
	return cache.NewCache("anidb_anime_by_id").TTL(8*time.Hour).Wrap(func() (any, error) {
		params := url.Values{}
		params.Set("aid", strconv.Itoa(animeID))
		return makeRequest("anime", params)
	})(ctx, animeID)
}

// GetHotAnime gets hot/trending anime from AniDB with 2-hour caching.
func GetHotAnime(ctx context.Context) (any, error) {
	return cache.NewCache("anidb_hot_anime").TTL(2 * time.Hour).Wrap(func() (any, error) {
		return makeRequest("hotanime", nil)
	})(ctx)
}

// GetRandomRecommendation gets a random anime recommendation with 4-hour caching.
func GetRandomRecommendation(ctx context.Context) (any, error) {
	return cache.NewCache("anidb_random_recommendation").TTL(4 * time.Hour).Wrap(func() (any, error) {
		return makeRequest("randomrecommendation", nil)
	})(ctx)
}

// GetRandomSimilar gets random similar anime pairs with 4-hour caching.
func GetRandomSimilar(ctx context.Context) (any, error) {
	return cache.NewCache("anidb_random_similar").TTL(4 * time.Hour).Wrap(func() (any, error) {
		return makeRequest("randomsimilar", nil)
	})(ctx)
}

// GetMainPageData gets combined main page data (hot, random recommendation, random similar) with 2-hour caching.
func GetMainPageData(ctx context.Context) (any, error) {
	return cache.NewCache("anidb_main").TTL(2 * time.Hour).Wrap(func() (any, error) {
		return makeRequest("main", nil)
	})(ctx)
}
