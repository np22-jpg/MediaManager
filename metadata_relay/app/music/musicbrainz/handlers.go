package musicbrainz

import (
	"log/slog"
	"net/http"
	"strconv"

	"relay/app/common"

	"github.com/gin-gonic/gin"
)

// handles MusicBrainz artist search route using Typesense
func SearchArtistsHandler(c *gin.Context) {
	slog.Debug("handling MusicBrainz artist search route")

	// Check if Typesense is available
	if typesenseClient == nil {
		common.WriteErrorResponse(c, "Search is not available - Typesense is not configured", http.StatusServiceUnavailable)
		return
	}

	query := c.Query("query")
	if query == "" {
		common.WriteErrorResponse(c, "query parameter is required", http.StatusBadRequest)
		return
	}

	limit := 25 // default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	result, err := SearchArtistsTypesense(c.Request.Context(), query, limit)
	if err != nil {
		slog.Error("failed to search artists", "error", err)
		common.WriteErrorResponse(c, "Failed to search artists", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(c, result)
}

// handles MusicBrainz get artist route using direct database lookup
func GetArtistHandler(c *gin.Context) {
	slog.Debug("handling MusicBrainz get artist route")

	mbid := c.Param("mbid")
	if mbid == "" {
		common.WriteErrorResponse(c, "mbid parameter is required", http.StatusBadRequest)
		return
	}

	result, err := GetArtist(c.Request.Context(), mbid)
	if err != nil {
		slog.Error("failed to get artist", "error", err)
		common.WriteErrorResponse(c, "Failed to fetch artist", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(c, result)
}

// handles MusicBrainz release group search route using Typesense
func SearchReleaseGroupsHandler(c *gin.Context) {
	slog.Debug("handling MusicBrainz release group search route")

	// Check if Typesense is available
	if typesenseClient == nil {
		common.WriteErrorResponse(c, "Search is not available - Typesense is not configured", http.StatusServiceUnavailable)
		return
	}

	query := c.Query("query")
	if query == "" {
		common.WriteErrorResponse(c, "query parameter is required", http.StatusBadRequest)
		return
	}

	limit := 25 // default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	result, err := SearchReleaseGroupsTypesense(c.Request.Context(), query, limit)
	if err != nil {
		slog.Error("failed to search release groups", "error", err)
		common.WriteErrorResponse(c, "Failed to search release groups", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(c, result)
}

// handles MusicBrainz get release group route
func GetReleaseGroupHandler(c *gin.Context) {
	slog.Debug("handling MusicBrainz get release group route")

	mbid := c.Param("mbid")
	if mbid == "" {
		common.WriteErrorResponse(c, "mbid parameter is required", http.StatusBadRequest)
		return
	}

	result, err := GetReleaseGroup(c.Request.Context(), mbid)
	if err != nil {
		slog.Error("failed to get release group", "error", err)
		common.WriteErrorResponse(c, "Failed to fetch release group", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(c, result)
}

// handles MusicBrainz release search route
func SearchReleasesHandler(c *gin.Context) {
	slog.Debug("handling MusicBrainz release search route")

	query := c.Query("query")
	if query == "" {
		common.WriteErrorResponse(c, "query parameter is required", http.StatusBadRequest)
		return
	}

	limit := 25 // default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Prefer Typesense search for releases
	if typesenseClient == nil {
		common.WriteErrorResponse(c, "Search is not available - Typesense is not configured", http.StatusServiceUnavailable)
		return
	}
	result, err := SearchReleasesTypesense(c.Request.Context(), query, limit)
	if err != nil {
		slog.Error("failed to search releases", "error", err)
		common.WriteErrorResponse(c, "Failed to search releases", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(c, result)
}

// handles MusicBrainz get release route
func GetReleaseHandler(c *gin.Context) {
	slog.Debug("handling MusicBrainz get release route")

	mbid := c.Param("mbid")
	if mbid == "" {
		common.WriteErrorResponse(c, "mbid parameter is required", http.StatusBadRequest)
		return
	}

	result, err := GetRelease(c.Request.Context(), mbid)
	if err != nil {
		slog.Error("failed to get release", "error", err)
		common.WriteErrorResponse(c, "Failed to fetch release", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(c, result)
}

// handles MusicBrainz recording search route using Typesense
func SearchRecordingsHandler(c *gin.Context) {
	slog.Debug("handling MusicBrainz recording search route")

	// Check if Typesense is available
	if typesenseClient == nil {
		common.WriteErrorResponse(c, "Search is not available - Typesense is not configured", http.StatusServiceUnavailable)
		return
	}

	query := c.Query("query")
	if query == "" {
		common.WriteErrorResponse(c, "query parameter is required", http.StatusBadRequest)
		return
	}

	limit := 25 // default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	result, err := SearchRecordingsTypesense(c.Request.Context(), query, limit)
	if err != nil {
		slog.Error("failed to search recordings", "error", err)
		common.WriteErrorResponse(c, "Failed to search recordings", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(c, result)
}

// handles MusicBrainz get recording route
func GetRecordingHandler(c *gin.Context) {
	slog.Debug("handling MusicBrainz get recording route")

	mbid := c.Param("mbid")
	if mbid == "" {
		common.WriteErrorResponse(c, "mbid parameter is required", http.StatusBadRequest)
		return
	}

	result, err := GetRecording(c.Request.Context(), mbid)
	if err != nil {
		slog.Error("failed to get recording", "error", err)
		common.WriteErrorResponse(c, "Failed to fetch recording", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(c, result)
}

// handles browsing release groups for an artist
func BrowseArtistReleaseGroupsHandler(c *gin.Context) {
	slog.Debug("handling MusicBrainz browse artist release groups route")

	artistMbid := c.Param("mbid")
	if artistMbid == "" {
		common.WriteErrorResponse(c, "artist mbid parameter is required", http.StatusBadRequest)
		return
	}

	limit := 25 // default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	result, err := BrowseArtistReleaseGroups(c.Request.Context(), artistMbid, limit)
	if err != nil {
		slog.Error("failed to browse artist release groups", "error", err)
		common.WriteErrorResponse(c, "Failed to fetch artist release groups", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(c, result)
}

// handles browsing releases for a release group
func BrowseReleaseGroupReleasesHandler(c *gin.Context) {
	slog.Debug("handling MusicBrainz browse release group releases route")

	releaseGroupMbid := c.Param("mbid")
	if releaseGroupMbid == "" {
		common.WriteErrorResponse(c, "release group mbid parameter is required", http.StatusBadRequest)
		return
	}

	limit := 25 // default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	result, err := BrowseReleaseGroupReleases(c.Request.Context(), releaseGroupMbid, limit)
	if err != nil {
		slog.Error("failed to browse release group releases", "error", err)
		common.WriteErrorResponse(c, "Failed to fetch release group releases", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(c, result)
}

// handles advanced artist search with multiple fields
func AdvancedSearchArtistsHandler(c *gin.Context) {
	slog.Debug("handling MusicBrainz advanced artist search route")

	artistName := c.Query("artist")
	area := c.Query("area")
	beginDate := c.Query("begin")
	endDate := c.Query("end")

	if artistName == "" && area == "" && beginDate == "" && endDate == "" {
		common.WriteErrorResponse(c, "at least one search parameter is required", http.StatusBadRequest)
		return
	}

	limit := 25 // default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	result, err := AdvancedSearchArtists(c.Request.Context(), artistName, area, beginDate, endDate, limit)
	if err != nil {
		slog.Error("failed to perform advanced artist search", "error", err)
		common.WriteErrorResponse(c, "Failed to perform advanced artist search", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(c, result)
}
