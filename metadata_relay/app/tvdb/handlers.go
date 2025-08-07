package tvdb

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

// handles TVDB trending TV route
func TrendingTVHandler(c *gin.Context) {
	slog.Debug("handling TVDB trending TV route")

	result, err := GetTrendingTV(c.Request.Context())
	if err != nil {
		slog.Error("failed to get trending TV", "error", err)
		writeErrorResponse(c, "Failed to fetch trending TV shows", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(c, result)
}

// handles TVDB search TV route
func SearchTVHandler(c *gin.Context) {
	slog.Debug("handling TVDB search TV route")

	query := c.Query("query")
	if query == "" {
		writeErrorResponse(c, "query parameter is required", http.StatusBadRequest)
		return
	}

	result, err := SearchTV(c.Request.Context(), query)
	if err != nil {
		slog.Error("failed to search TV", "error", err)
		writeErrorResponse(c, "Failed to search TV shows", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(c, result)
}

// handles TVDB get TV show route
func GetTVShowHandler(c *gin.Context) {
	slog.Debug("handling TVDB get TV show route")

	showID, err := strconv.Atoi(c.Param("showId"))
	if err != nil {
		writeErrorResponse(c, "Invalid show ID", http.StatusBadRequest)
		return
	}

	result, err := GetTVShow(c.Request.Context(), showID)
	if err != nil {
		slog.Error("failed to get TV show", "error", err)
		writeErrorResponse(c, "Failed to fetch TV show", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(c, result)
}

// handles TVDB get TV season route
func GetTVSeasonHandler(c *gin.Context) {
	slog.Debug("handling TVDB get TV season route")

	seasonID, err := strconv.Atoi(c.Param("seasonId"))
	if err != nil {
		writeErrorResponse(c, "Invalid season ID", http.StatusBadRequest)
		return
	}

	result, err := GetTVSeason(c.Request.Context(), seasonID)
	if err != nil {
		slog.Error("failed to get TV season", "error", err)
		writeErrorResponse(c, "Failed to fetch TV season", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(c, result)
}

// handles TVDB trending movies route
func TrendingMoviesHandler(c *gin.Context) {
	slog.Debug("handling TVDB trending movies route")

	result, err := GetTrendingMovies(c.Request.Context())
	if err != nil {
		slog.Error("failed to get trending movies", "error", err)
		writeErrorResponse(c, "Failed to fetch trending movies", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(c, result)
}

// handles TVDB search movies route
func SearchMoviesHandler(c *gin.Context) {
	slog.Debug("handling TVDB search movies route")

	query := c.Query("query")
	if query == "" {
		writeErrorResponse(c, "query parameter is required", http.StatusBadRequest)
		return
	}

	result, err := SearchMovies(c.Request.Context(), query)
	if err != nil {
		slog.Error("failed to search movies", "error", err)
		writeErrorResponse(c, "Failed to search movies", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(c, result)
}

// handles TVDB get movie route
func GetMovieHandler(c *gin.Context) {
	slog.Debug("handling TVDB get movie route")

	movieID, err := strconv.Atoi(c.Param("movieId"))
	if err != nil {
		writeErrorResponse(c, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	result, err := GetMovie(c.Request.Context(), movieID)
	if err != nil {
		slog.Error("failed to get movie", "error", err)
		writeErrorResponse(c, "Failed to fetch movie", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(c, result)
}
