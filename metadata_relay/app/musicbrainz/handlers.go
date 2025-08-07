package musicbrainz

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func writeJSONResponse(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}

func writeErrorResponse(c *gin.Context, message string, statusCode int) {
	c.JSON(statusCode, gin.H{"error": message})
}

// handles MusicBrainz artist search route
func SearchArtistsHandler(c *gin.Context) {
	slog.Debug("handling MusicBrainz artist search route")

	query := c.Query("query")
	if query == "" {
		writeErrorResponse(c, "query parameter is required", http.StatusBadRequest)
		return
	}

	limit := 25 // default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	result, err := SearchArtists(c.Request.Context(), query, limit)
	if err != nil {
		slog.Error("failed to search artists", "error", err)
		writeErrorResponse(c, "Failed to search artists", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(c, result)
}

// handles MusicBrainz get artist route
func GetArtistHandler(c *gin.Context) {
	slog.Debug("handling MusicBrainz get artist route")

	mbid := c.Param("mbid")
	if mbid == "" {
		writeErrorResponse(c, "mbid parameter is required", http.StatusBadRequest)
		return
	}

	result, err := GetArtist(c.Request.Context(), mbid)
	if err != nil {
		slog.Error("failed to get artist", "error", err)
		writeErrorResponse(c, "Failed to fetch artist", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(c, result)
}

// handles MusicBrainz release group search route
func SearchReleaseGroupsHandler(c *gin.Context) {
	slog.Debug("handling MusicBrainz release group search route")

	query := c.Query("query")
	if query == "" {
		writeErrorResponse(c, "query parameter is required", http.StatusBadRequest)
		return
	}

	limit := 25 // default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	result, err := SearchReleaseGroups(c.Request.Context(), query, limit)
	if err != nil {
		slog.Error("failed to search release groups", "error", err)
		writeErrorResponse(c, "Failed to search release groups", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(c, result)
}

// handles MusicBrainz get release group route
func GetReleaseGroupHandler(c *gin.Context) {
	slog.Debug("handling MusicBrainz get release group route")

	mbid := c.Param("mbid")
	if mbid == "" {
		writeErrorResponse(c, "mbid parameter is required", http.StatusBadRequest)
		return
	}

	result, err := GetReleaseGroup(c.Request.Context(), mbid)
	if err != nil {
		slog.Error("failed to get release group", "error", err)
		writeErrorResponse(c, "Failed to fetch release group", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(c, result)
}

// handles MusicBrainz release search route
func SearchReleasesHandler(c *gin.Context) {
	slog.Debug("handling MusicBrainz release search route")

	query := c.Query("query")
	if query == "" {
		writeErrorResponse(c, "query parameter is required", http.StatusBadRequest)
		return
	}

	limit := 25 // default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	result, err := SearchReleases(c.Request.Context(), query, limit)
	if err != nil {
		slog.Error("failed to search releases", "error", err)
		writeErrorResponse(c, "Failed to search releases", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(c, result)
}

// handles MusicBrainz get release route
func GetReleaseHandler(c *gin.Context) {
	slog.Debug("handling MusicBrainz get release route")

	mbid := c.Param("mbid")
	if mbid == "" {
		writeErrorResponse(c, "mbid parameter is required", http.StatusBadRequest)
		return
	}

	result, err := GetRelease(c.Request.Context(), mbid)
	if err != nil {
		slog.Error("failed to get release", "error", err)
		writeErrorResponse(c, "Failed to fetch release", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(c, result)
}

// handles MusicBrainz recording search route
func SearchRecordingsHandler(c *gin.Context) {
	slog.Debug("handling MusicBrainz recording search route")

	query := c.Query("query")
	if query == "" {
		writeErrorResponse(c, "query parameter is required", http.StatusBadRequest)
		return
	}

	limit := 25 // default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	result, err := SearchRecordings(c.Request.Context(), query, limit)
	if err != nil {
		slog.Error("failed to search recordings", "error", err)
		writeErrorResponse(c, "Failed to search recordings", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(c, result)
}

// handles MusicBrainz get recording route
func GetRecordingHandler(c *gin.Context) {
	slog.Debug("handling MusicBrainz get recording route")

	mbid := c.Param("mbid")
	if mbid == "" {
		writeErrorResponse(c, "mbid parameter is required", http.StatusBadRequest)
		return
	}

	result, err := GetRecording(c.Request.Context(), mbid)
	if err != nil {
		slog.Error("failed to get recording", "error", err)
		writeErrorResponse(c, "Failed to fetch recording", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(c, result)
}

// handles browsing release groups for an artist
func BrowseArtistReleaseGroupsHandler(c *gin.Context) {
	slog.Debug("handling MusicBrainz browse artist release groups route")

	artistMbid := c.Param("mbid")
	if artistMbid == "" {
		writeErrorResponse(c, "artist mbid parameter is required", http.StatusBadRequest)
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
		writeErrorResponse(c, "Failed to fetch artist release groups", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(c, result)
}

// handles browsing releases for a release group
func BrowseReleaseGroupReleasesHandler(c *gin.Context) {
	slog.Debug("handling MusicBrainz browse release group releases route")

	releaseGroupMbid := c.Param("mbid")
	if releaseGroupMbid == "" {
		writeErrorResponse(c, "release group mbid parameter is required", http.StatusBadRequest)
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
		writeErrorResponse(c, "Failed to fetch release group releases", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(c, result)
}

// handles advanced artist search with multiple fields
func AdvancedSearchArtistsHandler(c *gin.Context) {
	slog.Debug("handling MusicBrainz advanced artist search route")

	artistName := c.Query("artist")
	area := c.Query("area")
	beginDate := c.Query("begin")
	endDate := c.Query("end")

	if artistName == "" && area == "" && beginDate == "" && endDate == "" {
		writeErrorResponse(c, "at least one search parameter is required", http.StatusBadRequest)
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
		writeErrorResponse(c, "Failed to perform advanced artist search", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(c, result)
}
