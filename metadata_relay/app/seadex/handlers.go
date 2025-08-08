package seadex

import (
	"log/slog"
	"net/http"
	"strconv"

	"relay/app/common"

	"github.com/gin-gonic/gin"
)

// SearchEntriesHandler handles SeaDex anime entry search endpoint.
func SearchEntriesHandler(c *gin.Context) {
	slog.Debug("handling SeaDex search entries route")

	query := c.Query("query")
	if query == "" {
		common.WriteErrorResponse(c, "query parameter is required", http.StatusBadRequest)
		return
	}

	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	perPage := 30
	if perPageStr := c.Query("perPage"); perPageStr != "" {
		if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 && pp <= 100 {
			perPage = pp
		}
	}

	result, err := SearchEntries(c.Request.Context(), query, page, perPage)
	if err != nil {
		slog.Error("failed to search SeaDex entries", "error", err)
		common.WriteErrorResponse(c, "Failed to search anime entries", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(c, result)
}

// GetEntryByIDHandler handles SeaDex individual entry details endpoint.
func GetEntryByIDHandler(c *gin.Context) {
	slog.Debug("handling SeaDex get entry by ID route")

	entryID := c.Param("id")
	if entryID == "" {
		common.WriteErrorResponse(c, "Entry ID is required", http.StatusBadRequest)
		return
	}

	result, err := GetEntryByID(c.Request.Context(), entryID)
	if err != nil {
		slog.Error("failed to get SeaDex entry", "error", err, "id", entryID)
		common.WriteErrorResponse(c, "Failed to fetch anime entry", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(c, result)
}

// GetEntryByAnilistIDHandler handles SeaDex entry lookup by AniList ID.
func GetEntryByAnilistIDHandler(c *gin.Context) {
	slog.Debug("handling SeaDx get entry by AniList ID route")

	anilistIDStr := c.Param("anilistId")
	anilistID, err := strconv.Atoi(anilistIDStr)
	if err != nil {
		common.WriteErrorResponse(c, "Invalid AniList ID", http.StatusBadRequest)
		return
	}

	result, err := GetEntryByAnilistID(c.Request.Context(), anilistID)
	if err != nil {
		slog.Error("failed to get SeaDex entry by AniList ID", "error", err, "anilist_id", anilistID)
		common.WriteErrorResponse(c, "Failed to fetch anime entry", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(c, result)
}

// GetTrendingEntriesHandler handles SeaDex trending anime entries endpoint.
func GetTrendingEntriesHandler(c *gin.Context) {
	slog.Debug("handling SeaDex trending entries route")

	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	result, err := GetTrendingEntries(c.Request.Context(), limit)
	if err != nil {
		slog.Error("failed to get trending SeaDex entries", "error", err)
		common.WriteErrorResponse(c, "Failed to fetch trending anime entries", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(c, result)
}

// GetEntriesByReleaseGroupHandler handles SeaDex entries search by release group.
func GetEntriesByReleaseGroupHandler(c *gin.Context) {
	slog.Debug("handling SeaDex entries by release group route")

	releaseGroup := c.Query("group")
	if releaseGroup == "" {
		common.WriteErrorResponse(c, "release group parameter is required", http.StatusBadRequest)
		return
	}

	result, err := GetEntriesByReleaseGroup(c.Request.Context(), releaseGroup)
	if err != nil {
		slog.Error("failed to get SeaDex entries by release group", "error", err, "group", releaseGroup)
		common.WriteErrorResponse(c, "Failed to fetch entries by release group", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(c, result)
}

// GetEntriesByTrackerHandler handles SeaDex entries search by tracker.
func GetEntriesByTrackerHandler(c *gin.Context) {
	slog.Debug("handling SeaDex entries by tracker route")

	tracker := c.Query("tracker")
	if tracker == "" {
		common.WriteErrorResponse(c, "tracker parameter is required", http.StatusBadRequest)
		return
	}

	result, err := GetEntriesByTracker(c.Request.Context(), tracker)
	if err != nil {
		slog.Error("failed to get SeaDex entries by tracker", "error", err, "tracker", tracker)
		common.WriteErrorResponse(c, "Failed to fetch entries by tracker", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(c, result)
}
