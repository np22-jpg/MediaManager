package tmdb

import (
	"log/slog"
	"net/http"
	"strconv"

	"relay/app/common"

	"github.com/gin-gonic/gin"
)

// TrendingTVHandler handles TMDB trending TV shows endpoint.
func TrendingTVHandler(c *gin.Context) {
	slog.Debug("handling TMDB trending TV route")

	result, err := GetTrendingTV(c.Request.Context())
	if err != nil {
		slog.Error("failed to get trending TV", "error", err)
		common.WriteErrorResponse(c, "Failed to fetch trending TV shows", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(c, result)
}

// SearchTVHandler handles TMDB TV show search endpoint with query and pagination.
func SearchTVHandler(c *gin.Context) {
	slog.Debug("handling TMDB search TV route")

	query := c.Query("query")
	if query == "" {
		common.WriteErrorResponse(c, "query parameter is required", http.StatusBadRequest)
		return
	}

	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}

	result, err := SearchTV(c.Request.Context(), query, page)
	if err != nil {
		slog.Error("failed to search TV", "error", err)
		common.WriteErrorResponse(c, "Failed to search TV shows", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(c, result)
}

// GetTVShowHandler handles TMDB individual TV show details endpoint.
func GetTVShowHandler(c *gin.Context) {
	slog.Debug("handling TMDB get TV show route")

	showID, err := strconv.Atoi(c.Param("showId"))
	if err != nil {
		common.WriteErrorResponse(c, "Invalid show ID", http.StatusBadRequest)
		return
	}

	result, err := GetTVShow(c.Request.Context(), showID)
	if err != nil {
		slog.Error("failed to get TV show", "error", err)
		common.WriteErrorResponse(c, "Failed to fetch TV show", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(c, result)
}

// GetTVSeasonHandler handles TMDB TV season details endpoint.
func GetTVSeasonHandler(c *gin.Context) {
	slog.Debug("handling TMDB get TV season route")

	showID, err := strconv.Atoi(c.Param("showId"))
	if err != nil {
		common.WriteErrorResponse(c, "Invalid show ID", http.StatusBadRequest)
		return
	}

	seasonNumber, err := strconv.Atoi(c.Param("seasonNumber"))
	if err != nil {
		common.WriteErrorResponse(c, "Invalid season number", http.StatusBadRequest)
		return
	}

	result, err := GetTVSeason(c.Request.Context(), showID, seasonNumber)
	if err != nil {
		slog.Error("failed to get TV season", "error", err)
		common.WriteErrorResponse(c, "Failed to fetch TV season", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(c, result)
}

// TrendingMoviesHandler handles TMDB trending movies endpoint.
func TrendingMoviesHandler(c *gin.Context) {
	slog.Debug("handling TMDB trending movies route")

	result, err := GetTrendingMovies(c.Request.Context())
	if err != nil {
		slog.Error("failed to get trending movies", "error", err)
		common.WriteErrorResponse(c, "Failed to fetch trending movies", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(c, result)
}

// SearchMoviesHandler handles TMDB movie search endpoint with query and pagination.
func SearchMoviesHandler(c *gin.Context) {
	slog.Debug("handling TMDB search movies route")

	query := c.Query("query")
	if query == "" {
		common.WriteErrorResponse(c, "query parameter is required", http.StatusBadRequest)
		return
	}

	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}

	result, err := SearchMovies(c.Request.Context(), query, page)
	if err != nil {
		slog.Error("failed to search movies", "error", err)
		common.WriteErrorResponse(c, "Failed to search movies", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(c, result)
}

// GetMovieHandler handles TMDB individual movie details endpoint.
func GetMovieHandler(c *gin.Context) {
	slog.Debug("handling TMDB get movie route")

	movieID, err := strconv.Atoi(c.Param("movieId"))
	if err != nil {
		common.WriteErrorResponse(c, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	result, err := GetMovie(c.Request.Context(), movieID)
	if err != nil {
		slog.Error("failed to get movie", "error", err)
		common.WriteErrorResponse(c, "Failed to fetch movie", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(c, result)
}
