package anidb

import (
	"log/slog"
	"net/http"
	"strconv"

	"relay/app/common"

	"github.com/gin-gonic/gin"
)

// GetAnimeByIDHandler handles AniDB anime lookup by ID endpoint.
func GetAnimeByIDHandler(c *gin.Context) {
	slog.Debug("handling AniDB get anime by ID route")

	animeIDStr := c.Param("id")
	animeID, err := strconv.Atoi(animeIDStr)
	if err != nil {
		common.WriteErrorResponse(c, "Invalid anime ID", http.StatusBadRequest)
		return
	}

	result, err := GetAnimeByID(c.Request.Context(), animeID)
	if err != nil {
		slog.Error("failed to get AniDB anime", "error", err, "id", animeID)
		common.WriteErrorResponse(c, "Failed to fetch anime", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(c, result)
}

// GetHotAnimeHandler handles AniDB hot/trending anime endpoint.
func GetHotAnimeHandler(c *gin.Context) {
	slog.Debug("handling AniDB hot anime route")

	result, err := GetHotAnime(c.Request.Context())
	if err != nil {
		slog.Error("failed to get hot anime", "error", err)
		common.WriteErrorResponse(c, "Failed to fetch hot anime", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(c, result)
}

// GetRandomRecommendationHandler handles AniDB random recommendation endpoint.
func GetRandomRecommendationHandler(c *gin.Context) {
	slog.Debug("handling AniDB random recommendation route")

	result, err := GetRandomRecommendation(c.Request.Context())
	if err != nil {
		slog.Error("failed to get random recommendation", "error", err)
		common.WriteErrorResponse(c, "Failed to fetch recommendation", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(c, result)
}

// GetRandomSimilarHandler handles AniDB random similar anime endpoint.
func GetRandomSimilarHandler(c *gin.Context) {
	slog.Debug("handling AniDB random similar route")

	result, err := GetRandomSimilar(c.Request.Context())
	if err != nil {
		slog.Error("failed to get random similar", "error", err)
		common.WriteErrorResponse(c, "Failed to fetch similar anime", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(c, result)
}

// GetMainPageDataHandler handles AniDB main page data endpoint.
func GetMainPageDataHandler(c *gin.Context) {
	slog.Debug("handling AniDB main page data route")

	result, err := GetMainPageData(c.Request.Context())
	if err != nil {
		slog.Error("failed to get main page data", "error", err)
		common.WriteErrorResponse(c, "Failed to fetch main page data", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(c, result)
}
