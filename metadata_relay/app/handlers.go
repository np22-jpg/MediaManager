package app

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"relay/app/tmdb"
	"relay/app/tvdb"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	// Only respond to the exact root path
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	slog.Debug("handling hello route", "method", r.Method, "path", r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Hello World"}); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}

func writeJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": message}); err != nil {
		slog.Error("failed to encode error response", "error", err)
	}
}

// TMDB Handlers

func TMDBTrendingTVHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling TMDB trending TV route")

	result, err := tmdb.GetTrendingTV(r.Context())
	if err != nil {
		slog.Error("failed to get trending TV", "error", err)
		writeErrorResponse(w, "Failed to fetch trending TV shows", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, result)
}

func TMDBSearchTVHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling TMDB search TV route")

	query := r.URL.Query().Get("query")
	if query == "" {
		writeErrorResponse(w, "query parameter is required", http.StatusBadRequest)
		return
	}

	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}

	result, err := tmdb.SearchTV(r.Context(), query, page)
	if err != nil {
		slog.Error("failed to search TV", "error", err)
		writeErrorResponse(w, "Failed to search TV shows", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, result)
}

func TMDBGetTVShowHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling TMDB get TV show route")

	showIDStr := r.PathValue("showId")
	showID, err := strconv.Atoi(showIDStr)
	if err != nil {
		writeErrorResponse(w, "Invalid show ID", http.StatusBadRequest)
		return
	}

	result, err := tmdb.GetTVShow(r.Context(), showID)
	if err != nil {
		slog.Error("failed to get TV show", "error", err)
		writeErrorResponse(w, "Failed to fetch TV show", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, result)
}

func TMDBGetTVSeasonHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling TMDB get TV season route")

	showIDStr := r.PathValue("showId")
	seasonNumberStr := r.PathValue("seasonNumber")

	showID, err := strconv.Atoi(showIDStr)
	if err != nil {
		writeErrorResponse(w, "Invalid show ID", http.StatusBadRequest)
		return
	}

	seasonNumber, err := strconv.Atoi(seasonNumberStr)
	if err != nil {
		writeErrorResponse(w, "Invalid season number", http.StatusBadRequest)
		return
	}

	result, err := tmdb.GetTVSeason(r.Context(), showID, seasonNumber)
	if err != nil {
		slog.Error("failed to get TV season", "error", err)
		writeErrorResponse(w, "Failed to fetch TV season", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, result)
}

func TMDBTrendingMoviesHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling TMDB trending movies route")

	result, err := tmdb.GetTrendingMovies(r.Context())
	if err != nil {
		slog.Error("failed to get trending movies", "error", err)
		writeErrorResponse(w, "Failed to fetch trending movies", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, result)
}

func TMDBSearchMoviesHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling TMDB search movies route")

	query := r.URL.Query().Get("query")
	if query == "" {
		writeErrorResponse(w, "query parameter is required", http.StatusBadRequest)
		return
	}

	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}

	result, err := tmdb.SearchMovies(r.Context(), query, page)
	if err != nil {
		slog.Error("failed to search movies", "error", err)
		writeErrorResponse(w, "Failed to search movies", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, result)
}

func TMDBGetMovieHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling TMDB get movie route")

	movieIDStr := r.PathValue("movieId")
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		writeErrorResponse(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	result, err := tmdb.GetMovie(r.Context(), movieID)
	if err != nil {
		slog.Error("failed to get movie", "error", err)
		writeErrorResponse(w, "Failed to fetch movie", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, result)
}

// TVDB Handlers

func TVDBTrendingTVHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling TVDB trending TV route")

	result, err := tvdb.GetTrendingTV(r.Context())
	if err != nil {
		slog.Error("failed to get trending TV", "error", err)
		writeErrorResponse(w, "Failed to fetch trending TV shows", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, result)
}

func TVDBSearchTVHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling TVDB search TV route")

	query := r.URL.Query().Get("query")
	if query == "" {
		writeErrorResponse(w, "query parameter is required", http.StatusBadRequest)
		return
	}

	result, err := tvdb.SearchTV(r.Context(), query)
	if err != nil {
		slog.Error("failed to search TV", "error", err)
		writeErrorResponse(w, "Failed to search TV shows", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, result)
}

func TVDBGetTVShowHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling TVDB get TV show route")

	showIDStr := r.PathValue("showId")
	showID, err := strconv.Atoi(showIDStr)
	if err != nil {
		writeErrorResponse(w, "Invalid show ID", http.StatusBadRequest)
		return
	}

	result, err := tvdb.GetTVShow(r.Context(), showID)
	if err != nil {
		slog.Error("failed to get TV show", "error", err)
		writeErrorResponse(w, "Failed to fetch TV show", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, result)
}

func TVDBGetTVSeasonHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling TVDB get TV season route")

	seasonIDStr := r.PathValue("seasonId")
	seasonID, err := strconv.Atoi(seasonIDStr)
	if err != nil {
		writeErrorResponse(w, "Invalid season ID", http.StatusBadRequest)
		return
	}

	result, err := tvdb.GetTVSeason(r.Context(), seasonID)
	if err != nil {
		slog.Error("failed to get TV season", "error", err)
		writeErrorResponse(w, "Failed to fetch TV season", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, result)
}

func TVDBTrendingMoviesHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling TVDB trending movies route")

	result, err := tvdb.GetTrendingMovies(r.Context())
	if err != nil {
		slog.Error("failed to get trending movies", "error", err)
		writeErrorResponse(w, "Failed to fetch trending movies", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, result)
}

func TVDBSearchMoviesHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling TVDB search movies route")

	query := r.URL.Query().Get("query")
	if query == "" {
		writeErrorResponse(w, "query parameter is required", http.StatusBadRequest)
		return
	}

	result, err := tvdb.SearchMovies(r.Context(), query)
	if err != nil {
		slog.Error("failed to search movies", "error", err)
		writeErrorResponse(w, "Failed to search movies", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, result)
}

func TVDBGetMovieHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling TVDB get movie route")

	movieIDStr := r.PathValue("movieId")
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		writeErrorResponse(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	result, err := tvdb.GetMovie(r.Context(), movieID)
	if err != nil {
		slog.Error("failed to get movie", "error", err)
		writeErrorResponse(w, "Failed to fetch movie", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, result)
}
