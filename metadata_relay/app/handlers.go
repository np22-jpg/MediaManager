package app

import (
	"log/slog"
	"net/http"
	"strconv"

	"relay/app/tmdb"
	"relay/app/tvdb"

	"github.com/gin-gonic/gin"
)

func HelloHandler(c *gin.Context) {
	slog.Debug("handling hello route", "method", c.Request.Method, "path", c.Request.URL.Path)
	c.JSON(http.StatusOK, gin.H{"message": "Hello World"})
}

func writeJSONResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

func writeErrorResponse(c *gin.Context, message string, statusCode int) {
	c.JSON(statusCode, gin.H{"error": message})
}

// TMDB Handlers

func TMDBTrendingTVHandler(c *gin.Context) {
	slog.Debug("handling TMDB trending TV route")

	result, err := tmdb.GetTrendingTV(c.Request.Context())
	if err != nil {
		slog.Error("failed to get trending TV", "error", err)
		writeErrorResponse(c, "Failed to fetch trending TV shows", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(c, result)
}

func TMDBSearchTVHandler(c *gin.Context) {
	slog.Debug("handling TMDB search TV route")

	query := c.Query("query")
	if query == "" {
		writeErrorResponse(c, "query parameter is required", http.StatusBadRequest)
		return
	}

	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}

	result, err := tmdb.SearchTV(c.Request.Context(), query, page)
	if err != nil {
		slog.Error("failed to search TV", "error", err)
		writeErrorResponse(c, "Failed to search TV shows", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(c, result)
}

func TMDBGetTVShowHandler(c *gin.Context) {
	slog.Debug("handling TMDB get TV show route")

	showID, err := strconv.Atoi(c.Param("showId"))
	if err != nil {
		writeErrorResponse(c, "Invalid show ID", http.StatusBadRequest)
		return
	}

	result, err := tmdb.GetTVShow(c.Request.Context(), showID)
	if err != nil {
		slog.Error("failed to get TV show", "error", err)
		writeErrorResponse(c, "Failed to fetch TV show", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(c, result)
}

func TMDBGetTVSeasonHandler(c *gin.Context) {
	slog.Debug("handling TMDB get TV season route")

	showID, err := strconv.Atoi(c.Param("showId"))
	if err != nil {
		writeErrorResponse(c, "Invalid show ID", http.StatusBadRequest)
		return
	}

	seasonNumber, err := strconv.Atoi(c.Param("seasonNumber"))
	if err != nil {
		writeErrorResponse(c, "Invalid season number", http.StatusBadRequest)
		return
	}

	result, err := tmdb.GetTVSeason(c.Request.Context(), showID, seasonNumber)
	if err != nil {
		slog.Error("failed to get TV season", "error", err)
		writeErrorResponse(c, "Failed to fetch TV season", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(c, result)
}

func TMDBTrendingMoviesHandler(c *gin.Context) {
	slog.Debug("handling TMDB trending movies route")

	result, err := tmdb.GetTrendingMovies(c.Request.Context())
	if err != nil {
		slog.Error("failed to get trending movies", "error", err)
		writeErrorResponse(c, "Failed to fetch trending movies", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(c, result)
}

func TMDBSearchMoviesHandler(c *gin.Context) {
	slog.Debug("handling TMDB search movies route")

	query := c.Query("query")
	if query == "" {
		writeErrorResponse(c, "query parameter is required", http.StatusBadRequest)
		return
	}

	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}

	result, err := tmdb.SearchMovies(c.Request.Context(), query, page)
	if err != nil {
		slog.Error("failed to search movies", "error", err)
		writeErrorResponse(c, "Failed to search movies", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(c, result)
}

func TMDBGetMovieHandler(c *gin.Context) {
	slog.Debug("handling TMDB get movie route")

	movieID, err := strconv.Atoi(c.Param("movieId"))
	if err != nil {
		writeErrorResponse(c, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	result, err := tmdb.GetMovie(c.Request.Context(), movieID)
	if err != nil {
		slog.Error("failed to get movie", "error", err)
		writeErrorResponse(c, "Failed to fetch movie", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(c, result)
}

// TVDB Handlers

func TVDBTrendingTVHandler(c *gin.Context) {
	slog.Debug("handling TVDB trending TV route")

	result, err := tvdb.GetTrendingTV(c.Request.Context())
	if err != nil {
		slog.Error("failed to get trending TV", "error", err)
		writeErrorResponse(c, "Failed to fetch trending TV shows", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(c, result)
}

func TVDBSearchTVHandler(c *gin.Context) {
	slog.Debug("handling TVDB search TV route")

	query := c.Query("query")
	if query == "" {
		writeErrorResponse(c, "query parameter is required", http.StatusBadRequest)
		return
	}

	result, err := tvdb.SearchTV(c.Request.Context(), query)
	if err != nil {
		slog.Error("failed to search TV", "error", err)
		writeErrorResponse(c, "Failed to search TV shows", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(c, result)
}

func TVDBGetTVShowHandler(c *gin.Context) {
	slog.Debug("handling TVDB get TV show route")

	showID, err := strconv.Atoi(c.Param("showId"))
	if err != nil {
		writeErrorResponse(c, "Invalid show ID", http.StatusBadRequest)
		return
	}

	result, err := tvdb.GetTVShow(c.Request.Context(), showID)
	if err != nil {
		slog.Error("failed to get TV show", "error", err)
		writeErrorResponse(c, "Failed to fetch TV show", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(c, result)
}

func TVDBGetTVSeasonHandler(c *gin.Context) {
	slog.Debug("handling TVDB get TV season route")

	seasonID, err := strconv.Atoi(c.Param("seasonId"))
	if err != nil {
		writeErrorResponse(c, "Invalid season ID", http.StatusBadRequest)
		return
	}

	result, err := tvdb.GetTVSeason(c.Request.Context(), seasonID)
	if err != nil {
		slog.Error("failed to get TV season", "error", err)
		writeErrorResponse(c, "Failed to fetch TV season", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(c, result)
}

func TVDBTrendingMoviesHandler(c *gin.Context) {
	slog.Debug("handling TVDB trending movies route")

	result, err := tvdb.GetTrendingMovies(c.Request.Context())
	if err != nil {
		slog.Error("failed to get trending movies", "error", err)
		writeErrorResponse(c, "Failed to fetch trending movies", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(c, result)
}

func TVDBSearchMoviesHandler(c *gin.Context) {
	slog.Debug("handling TVDB search movies route")

	query := c.Query("query")
	if query == "" {
		writeErrorResponse(c, "query parameter is required", http.StatusBadRequest)
		return
	}

	result, err := tvdb.SearchMovies(c.Request.Context(), query)
	if err != nil {
		slog.Error("failed to search movies", "error", err)
		writeErrorResponse(c, "Failed to search movies", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(c, result)
}

func TVDBGetMovieHandler(c *gin.Context) {
	slog.Debug("handling TVDB get movie route")

	movieID, err := strconv.Atoi(c.Param("movieId"))
	if err != nil {
		writeErrorResponse(c, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	result, err := tvdb.GetMovie(c.Request.Context(), movieID)
	if err != nil {
		slog.Error("failed to get movie", "error", err)
		writeErrorResponse(c, "Failed to fetch movie", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(c, result)
}
