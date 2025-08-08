package seadex

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

// SearchEntriesHandler handles SeaDex anime entry search endpoint.
func SearchEntriesHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling SeaDex search entries route")

	query := r.URL.Query().Get("query")
	if query == "" {
		common.WriteErrorResponse(w, "query parameter is required", http.StatusBadRequest)
		return
	}

	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	perPage := 30
	if perPageStr := r.URL.Query().Get("perPage"); perPageStr != "" {
		if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 && pp <= 100 {
			perPage = pp
		}
	}

	result, err := SearchEntries(r.Context(), query, page, perPage)
	if err != nil {
		slog.Error("failed to search SeaDex entries", "error", err)
		common.WriteErrorResponse(w, "Failed to search anime entries", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}

// GetEntryByIDHandler handles SeaDex individual entry details endpoint.
func GetEntryByIDHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling SeaDex get entry by ID route")

	entryID := GetURLParam(r, "id")
	if entryID == "" {
		common.WriteErrorResponse(w, "Entry ID is required", http.StatusBadRequest)
		return
	}

	result, err := GetEntryByID(r.Context(), entryID)
	if err != nil {
		slog.Error("failed to get SeaDex entry", "error", err, "id", entryID)
		common.WriteErrorResponse(w, "Failed to fetch anime entry", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}

// GetEntryByAnilistIDHandler handles SeaDex entry lookup by AniList ID.
func GetEntryByAnilistIDHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling SeaDx get entry by AniList ID route")

	anilistIDStr := GetURLParam(r, "anilistId")
	anilistID, err := strconv.Atoi(anilistIDStr)
	if err != nil {
		common.WriteErrorResponse(w, "Invalid AniList ID", http.StatusBadRequest)
		return
	}

	result, err := GetEntryByAnilistID(r.Context(), anilistID)
	if err != nil {
		slog.Error("failed to get SeaDex entry by AniList ID", "error", err, "anilist_id", anilistID)
		common.WriteErrorResponse(w, "Failed to fetch anime entry", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}

// GetTrendingEntriesHandler handles SeaDex trending anime entries endpoint.
func GetTrendingEntriesHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling SeaDex trending entries route")

	limit := 50
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	result, err := GetTrendingEntries(r.Context(), limit)
	if err != nil {
		slog.Error("failed to get trending SeaDex entries", "error", err)
		common.WriteErrorResponse(w, "Failed to fetch trending anime entries", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}

// GetEntriesByReleaseGroupHandler handles SeaDex entries search by release group.
func GetEntriesByReleaseGroupHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling SeaDex entries by release group route")

	releaseGroup := r.URL.Query().Get("group")
	if releaseGroup == "" {
		common.WriteErrorResponse(w, "release group parameter is required", http.StatusBadRequest)
		return
	}

	result, err := GetEntriesByReleaseGroup(r.Context(), releaseGroup)
	if err != nil {
		slog.Error("failed to get SeaDex entries by release group", "error", err, "group", releaseGroup)
		common.WriteErrorResponse(w, "Failed to fetch entries by release group", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}

// GetEntriesByTrackerHandler handles SeaDex entries search by tracker.
func GetEntriesByTrackerHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling SeaDex entries by tracker route")

	tracker := r.URL.Query().Get("tracker")
	if tracker == "" {
		common.WriteErrorResponse(w, "tracker parameter is required", http.StatusBadRequest)
		return
	}

	result, err := GetEntriesByTracker(r.Context(), tracker)
	if err != nil {
		slog.Error("failed to get SeaDex entries by tracker", "error", err, "tracker", tracker)
		common.WriteErrorResponse(w, "Failed to fetch entries by tracker", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}
