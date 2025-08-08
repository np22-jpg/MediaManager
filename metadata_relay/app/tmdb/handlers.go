package tmdb

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

// TrendingTVHandler handles TMDB trending TV shows endpoint.
func TrendingTVHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling TMDB trending TV route")

	result, err := GetTrendingTV(r.Context())
	if err != nil {
		slog.Error("failed to get trending TV", "error", err)
		common.WriteErrorResponse(w, "Failed to fetch trending TV shows", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}

// SearchTVHandler handles TMDB TV show search endpoint with query and pagination.
func SearchTVHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling TMDB search TV route")

	query := r.URL.Query().Get("query")
	if query == "" {
		common.WriteErrorResponse(w, "query parameter is required", http.StatusBadRequest)
		return
	}

	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}

	result, err := SearchTV(r.Context(), query, page)
	if err != nil {
		slog.Error("failed to search TV", "error", err)
		common.WriteErrorResponse(w, "Failed to search TV shows", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}

// GetTVShowHandler handles TMDB individual TV show details endpoint.
func GetTVShowHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling TMDB get TV show route")

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

// GetTVSeasonHandler handles TMDB TV season details endpoint.
func GetTVSeasonHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling TMDB get TV season route")

	showID, err := strconv.Atoi(GetURLParam(r, "showId"))
	if err != nil {
		common.WriteErrorResponse(w, "Invalid show ID", http.StatusBadRequest)
		return
	}

	seasonNumber, err := strconv.Atoi(GetURLParam(r, "seasonNumber"))
	if err != nil {
		common.WriteErrorResponse(w, "Invalid season number", http.StatusBadRequest)
		return
	}

	result, err := GetTVSeason(r.Context(), showID, seasonNumber)
	if err != nil {
		slog.Error("failed to get TV season", "error", err)
		common.WriteErrorResponse(w, "Failed to fetch TV season", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}

// TrendingMoviesHandler handles TMDB trending movies endpoint.
func TrendingMoviesHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling TMDB trending movies route")

	result, err := GetTrendingMovies(r.Context())
	if err != nil {
		slog.Error("failed to get trending movies", "error", err)
		common.WriteErrorResponse(w, "Failed to fetch trending movies", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}

// SearchMoviesHandler handles TMDB movie search endpoint with query and pagination.
func SearchMoviesHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling TMDB search movies route")

	query := r.URL.Query().Get("query")
	if query == "" {
		common.WriteErrorResponse(w, "query parameter is required", http.StatusBadRequest)
		return
	}

	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}

	result, err := SearchMovies(r.Context(), query, page)
	if err != nil {
		slog.Error("failed to search movies", "error", err)
		common.WriteErrorResponse(w, "Failed to search movies", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(w, result)
}

// GetMovieHandler handles TMDB individual movie details endpoint.
func GetMovieHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling TMDB get movie route")

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
