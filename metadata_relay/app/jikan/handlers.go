package jikan

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"relay/app/common"
)

// GetURLParam extracts a URL parameter from the request path
func GetURLParam(r *http.Request, key string) string {
	path := strings.TrimPrefix(r.URL.Path, "/jikan/")
	parts := strings.Split(path, "/")

	switch key {
	case "id":
		if len(parts) >= 2 {
			return parts[1]
		}
	}
	return ""
}

// GetAnimeByIDHandler handles Jikan anime lookup by ID endpoint.
func GetAnimeByIDHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling Jikan get anime by ID route")

	animeIDStr := GetURLParam(r, "id")
	animeID, err := strconv.Atoi(animeIDStr)
	if err != nil {
		common.WriteErrorResponse(w, "Invalid anime ID", http.StatusBadRequest)
		return
	}

	result, err := GetAnimeByID(r.Context(), animeID)
	if err != nil {
		slog.Error("failed to get Jikan anime", "error", err, "id", animeID)
		common.WriteErrorResponse(w, "Failed to fetch anime", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}

// GetTopAnimeHandler handles Jikan top anime endpoint.
func GetTopAnimeHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling Jikan top anime route")

	result, err := GetTopAnime(r.Context())
	if err != nil {
		slog.Error("failed to get top anime", "error", err)
		common.WriteErrorResponse(w, "Failed to fetch top anime", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}

// GetSeasonalAnimeHandler handles Jikan seasonal anime endpoint.
func GetSeasonalAnimeHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling Jikan seasonal anime route")

	result, err := GetSeasonalAnime(r.Context())
	if err != nil {
		slog.Error("failed to get seasonal anime", "error", err)
		common.WriteErrorResponse(w, "Failed to fetch seasonal anime", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}

// SearchAnimeHandler handles Jikan anime search endpoint.
func SearchAnimeHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling Jikan anime search route")

	query := r.URL.Query().Get("q")
	if query == "" {
		common.WriteErrorResponse(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	result, err := SearchAnime(r.Context(), query, page)
	if err != nil {
		slog.Error("failed to search anime", "error", err, "query", query)
		common.WriteErrorResponse(w, "Failed to search anime", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}

// GetAnimeRecommendationsHandler handles Jikan anime recommendations endpoint.
func GetAnimeRecommendationsHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling Jikan anime recommendations route")

	animeIDStr := GetURLParam(r, "id")
	animeID, err := strconv.Atoi(animeIDStr)
	if err != nil {
		common.WriteErrorResponse(w, "Invalid anime ID", http.StatusBadRequest)
		return
	}

	result, err := GetAnimeRecommendations(r.Context(), animeID)
	if err != nil {
		slog.Error("failed to get anime recommendations", "error", err, "id", animeID)
		common.WriteErrorResponse(w, "Failed to fetch recommendations", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}

// GetRandomAnimeHandler handles Jikan random anime endpoint.
func GetRandomAnimeHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling Jikan random anime route")

	result, err := GetRandomAnime(r.Context())
	if err != nil {
		slog.Error("failed to get random anime", "error", err)
		common.WriteErrorResponse(w, "Failed to fetch random anime", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}
