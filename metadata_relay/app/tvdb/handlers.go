package tvdb

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

// TrendingTVHandler handles TVDB trending TV shows endpoint.
func TrendingTVHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling TVDB trending TV route")

	result, err := GetTrendingTV(r.Context())
	if err != nil {
		slog.Error("failed to get trending TV", "error", err)
		common.WriteErrorResponse(w, "Failed to fetch trending TV shows", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}

// handles TVDB search TV route
func SearchTVHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling TVDB search TV route")

	query := r.URL.Query().Get("query")
	if query == "" {
		common.WriteErrorResponse(w, "query parameter is required", http.StatusBadRequest)
		return
	}

	result, err := SearchTV(r.Context(), query)
	if err != nil {
		slog.Error("failed to search TV", "error", err)
		common.WriteErrorResponse(w, "Failed to search TV shows", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}

// handles TVDB get TV show route
func GetTVShowHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling TVDB get TV show route")

	showID, err := strconv.Atoi(GetURLParam(r, "showId"))
	if err != nil {
		common.WriteErrorResponse(w, "Invalid show ID", http.StatusBadRequest)
		return
	}

	result, err := GetTVShow(r.Context(), showID)
	if err != nil {
		slog.Error("failed to get TV show", "error", err)
		common.WriteErrorResponse(w, "Failed to fetch TV show", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}

// handles TVDB get TV season route
func GetTVSeasonHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling TVDB get TV season route")

	seasonID, err := strconv.Atoi(GetURLParam(r, "seasonId"))
	if err != nil {
		common.WriteErrorResponse(w, "Invalid season ID", http.StatusBadRequest)
		return
	}

	result, err := GetTVSeason(r.Context(), seasonID)
	if err != nil {
		slog.Error("failed to get TV season", "error", err)
		common.WriteErrorResponse(w, "Failed to fetch TV season", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}

// handles TVDB trending movies route
func TrendingMoviesHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling TVDB trending movies route")

	result, err := GetTrendingMovies(r.Context())
	if err != nil {
		slog.Error("failed to get trending movies", "error", err)
		common.WriteErrorResponse(w, "Failed to fetch trending movies", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}

// handles TVDB search movies route
func SearchMoviesHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling TVDB search movies route")

	query := r.URL.Query().Get("query")
	if query == "" {
		common.WriteErrorResponse(w, "query parameter is required", http.StatusBadRequest)
		return
	}

	result, err := SearchMovies(r.Context(), query)
	if err != nil {
		slog.Error("failed to search movies", "error", err)
		common.WriteErrorResponse(w, "Failed to search movies", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}

// handles TVDB get movie route
func GetMovieHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling TVDB get movie route")

	movieID, err := strconv.Atoi(GetURLParam(r, "movieId"))
	if err != nil {
		common.WriteErrorResponse(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	result, err := GetMovie(r.Context(), movieID)
	if err != nil {
		slog.Error("failed to get movie", "error", err)
		common.WriteErrorResponse(w, "Failed to fetch movie", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}
