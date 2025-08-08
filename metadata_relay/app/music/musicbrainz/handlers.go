package musicbrainz

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

// handles MusicBrainz artist search route using Typesense
func SearchArtistsHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling MusicBrainz artist search route")

	// Check if Typesense is available
	if typesenseClient == nil {
		common.WriteErrorResponse(w, "Search is not available - Typesense is not configured", http.StatusServiceUnavailable)
		return
	}

	query := r.URL.Query().Get("query")
	if query == "" {
		common.WriteErrorResponse(w, "query parameter is required", http.StatusBadRequest)
		return
	}

	limit := 25 // default limit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	result, err := SearchArtistsTypesense(r.Context(), query, limit)
	if err != nil {
		slog.Error("failed to search artists", "error", err)
		common.WriteErrorResponse(w, "Failed to search artists", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}

// handles MusicBrainz get artist route using direct database lookup
func GetArtistHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling MusicBrainz get artist route")

	mbid := GetURLParam(r, "mbid")
	if mbid == "" {
		common.WriteErrorResponse(w, "mbid parameter is required", http.StatusBadRequest)
		return
	}

	result, err := GetArtist(r.Context(), mbid)
	if err != nil {
		slog.Error("failed to get artist", "error", err)
		common.WriteErrorResponse(w, "Failed to fetch artist", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}

// handles MusicBrainz release group search route using Typesense
func SearchReleaseGroupsHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling MusicBrainz release group search route")

	// Check if Typesense is available
	if typesenseClient == nil {
		common.WriteErrorResponse(w, "Search is not available - Typesense is not configured", http.StatusServiceUnavailable)
		return
	}

	query := r.URL.Query().Get("query")
	if query == "" {
		common.WriteErrorResponse(w, "query parameter is required", http.StatusBadRequest)
		return
	}

	limit := 25 // default limit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	result, err := SearchReleaseGroupsTypesense(r.Context(), query, limit)
	if err != nil {
		slog.Error("failed to search release groups", "error", err)
		common.WriteErrorResponse(w, "Failed to search release groups", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}

// handles MusicBrainz get release group route
func GetReleaseGroupHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling MusicBrainz get release group route")

	mbid := GetURLParam(r, "mbid")
	if mbid == "" {
		common.WriteErrorResponse(w, "mbid parameter is required", http.StatusBadRequest)
		return
	}

	result, err := GetReleaseGroup(r.Context(), mbid)
	if err != nil {
		slog.Error("failed to get release group", "error", err)
		common.WriteErrorResponse(w, "Failed to fetch release group", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}

// handles MusicBrainz release search route
func SearchReleasesHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling MusicBrainz release search route")

	query := r.URL.Query().Get("query")
	if query == "" {
		common.WriteErrorResponse(w, "query parameter is required", http.StatusBadRequest)
		return
	}

	limit := 25 // default limit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Prefer Typesense search for releases
	if typesenseClient == nil {
		common.WriteErrorResponse(w, "Search is not available - Typesense is not configured", http.StatusServiceUnavailable)
		return
	}
	result, err := SearchReleasesTypesense(r.Context(), query, limit)
	if err != nil {
		slog.Error("failed to search releases", "error", err)
		common.WriteErrorResponse(w, "Failed to search releases", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}

// handles MusicBrainz get release route
func GetReleaseHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling MusicBrainz get release route")

	mbid := GetURLParam(r, "mbid")
	if mbid == "" {
		common.WriteErrorResponse(w, "mbid parameter is required", http.StatusBadRequest)
		return
	}

	result, err := GetRelease(r.Context(), mbid)
	if err != nil {
		slog.Error("failed to get release", "error", err)
		common.WriteErrorResponse(w, "Failed to fetch release", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}

// handles MusicBrainz recording search route using Typesense
func SearchRecordingsHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling MusicBrainz recording search route")

	// Check if Typesense is available
	if typesenseClient == nil {
		common.WriteErrorResponse(w, "Search is not available - Typesense is not configured", http.StatusServiceUnavailable)
		return
	}

	query := r.URL.Query().Get("query")
	if query == "" {
		common.WriteErrorResponse(w, "query parameter is required", http.StatusBadRequest)
		return
	}

	limit := 25 // default limit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	result, err := SearchRecordingsTypesense(r.Context(), query, limit)
	if err != nil {
		slog.Error("failed to search recordings", "error", err)
		common.WriteErrorResponse(w, "Failed to search recordings", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}

// handles MusicBrainz get recording route
func GetRecordingHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling MusicBrainz get recording route")

	mbid := GetURLParam(r, "mbid")
	if mbid == "" {
		common.WriteErrorResponse(w, "mbid parameter is required", http.StatusBadRequest)
		return
	}

	result, err := GetRecording(r.Context(), mbid)
	if err != nil {
		slog.Error("failed to get recording", "error", err)
		common.WriteErrorResponse(w, "Failed to fetch recording", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}

// handles browsing release groups for an artist
func BrowseArtistReleaseGroupsHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling MusicBrainz browse artist release groups route")

	artistMbid := GetURLParam(r, "mbid")
	if artistMbid == "" {
		common.WriteErrorResponse(w, "artist mbid parameter is required", http.StatusBadRequest)
		return
	}

	limit := 25 // default limit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	result, err := BrowseArtistReleaseGroups(r.Context(), artistMbid, limit)
	if err != nil {
		slog.Error("failed to browse artist release groups", "error", err)
		common.WriteErrorResponse(w, "Failed to fetch artist release groups", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}

// handles browsing releases for a release group
func BrowseReleaseGroupReleasesHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling MusicBrainz browse release group releases route")

	releaseGroupMbid := GetURLParam(r, "mbid")
	if releaseGroupMbid == "" {
		common.WriteErrorResponse(w, "release group mbid parameter is required", http.StatusBadRequest)
		return
	}

	limit := 25 // default limit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	result, err := BrowseReleaseGroupReleases(r.Context(), releaseGroupMbid, limit)
	if err != nil {
		slog.Error("failed to browse release group releases", "error", err)
		common.WriteErrorResponse(w, "Failed to fetch release group releases", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}

// handles advanced artist search with multiple fields
func AdvancedSearchArtistsHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling MusicBrainz advanced artist search route")

	artistName := r.URL.Query().Get("artist")
	area := r.URL.Query().Get("area")
	beginDate := r.URL.Query().Get("begin")
	endDate := r.URL.Query().Get("end")

	if artistName == "" && area == "" && beginDate == "" && endDate == "" {
		common.WriteErrorResponse(w, "at least one search parameter is required", http.StatusBadRequest)
		return
	}

	limit := 25 // default limit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	result, err := AdvancedSearchArtists(r.Context(), artistName, area, beginDate, endDate, limit)
	if err != nil {
		slog.Error("failed to perform advanced artist search", "error", err)
		common.WriteErrorResponse(w, "Failed to perform advanced artist search", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}
