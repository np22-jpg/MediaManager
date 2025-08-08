package anidb

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"relay/app/common"
)

// GetURLParam extracts a URL parameter from the request path
func GetURLParam(r *http.Request, key string) string {
	// Get path parameters from context if they exist
	if params := r.Context().Value("pathParams"); params != nil {
		if paramMap, ok := params.(map[string]string); ok {
			return paramMap[key]
		}
	}

	// Fallback: extract from URL path - simple approach for IDs at the end
	path := r.URL.Path
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

// GetAnimeByIDHandler handles AniDB anime lookup by ID endpoint.
func GetAnimeByIDHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling AniDB get anime by ID route")

	animeIDStr := GetURLParam(r, "id")
	animeID, err := strconv.Atoi(animeIDStr)
	if err != nil {
		common.WriteErrorResponse(w, "Invalid anime ID", http.StatusBadRequest)
		return
	}

	result, err := GetAnimeByID(r.Context(), animeID)
	if err != nil {
		slog.Error("failed to get AniDB anime", "error", err, "id", animeID)
		common.WriteErrorResponse(w, "Failed to fetch anime", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}

// GetHotAnimeHandler handles AniDB hot/trending anime endpoint.
func GetHotAnimeHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling AniDB hot anime route")

	result, err := GetHotAnime(r.Context())
	if err != nil {
		slog.Error("failed to get hot anime", "error", err)
		common.WriteErrorResponse(w, "Failed to fetch hot anime", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}

// GetRandomRecommendationHandler handles AniDB random recommendation endpoint.
func GetRandomRecommendationHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling AniDB random recommendation route")

	result, err := GetRandomRecommendation(r.Context())
	if err != nil {
		slog.Error("failed to get random recommendation", "error", err)
		common.WriteErrorResponse(w, "Failed to fetch recommendation", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}

// GetRandomSimilarHandler handles AniDB random similar anime endpoint.
func GetRandomSimilarHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling AniDB random similar route")

	result, err := GetRandomSimilar(r.Context())
	if err != nil {
		slog.Error("failed to get random similar", "error", err)
		common.WriteErrorResponse(w, "Failed to fetch similar anime", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}

// GetMainPageDataHandler handles AniDB main page data endpoint.
func GetMainPageDataHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling AniDB main page data route")

	result, err := GetMainPageData(r.Context())
	if err != nil {
		slog.Error("failed to get main page data", "error", err)
		common.WriteErrorResponse(w, "Failed to fetch main page data", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}
