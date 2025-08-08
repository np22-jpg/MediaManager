package tvdb

import (
	"log/slog"
	"net/http"
	"strconv"

	"relay/app/common"

	"github.com/gin-gonic/gin"
)

// TrendingTVHandler handles TVDB trending TV shows endpoint.
func TrendingTVHandler(c *gin.Context) {
	slog.Debug("handling TVDB trending TV route")

	result, err := GetTrendingTV(c.Request.Context())
	if err != nil {
		slog.Error("failed to get trending TV", "error", err)
		common.WriteErrorResponse(c, "Failed to fetch trending TV shows", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(c, result)
}

// handles TVDB search TV route
func SearchTVHandler(c *gin.Context) {
	slog.Debug("handling TVDB search TV route")

	query := c.Query("query")
	if query == "" {
		common.WriteErrorResponse(c, "query parameter is required", http.StatusBadRequest)
		return
	}

	result, err := SearchTV(c.Request.Context(), query)
	if err != nil {
		slog.Error("failed to search TV", "error", err)
		common.WriteErrorResponse(c, "Failed to search TV shows", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(c, result)
}

// handles TVDB get TV show route
func GetTVShowHandler(c *gin.Context) {
	slog.Debug("handling TVDB get TV show route")

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

// handles TVDB get TV season route
func GetTVSeasonHandler(c *gin.Context) {
	slog.Debug("handling TVDB get TV season route")

	seasonID, err := strconv.Atoi(c.Param("seasonId"))
	if err != nil {
		common.WriteErrorResponse(c, "Invalid season ID", http.StatusBadRequest)
		return
	}

	result, err := GetTVSeason(c.Request.Context(), seasonID)
	if err != nil {
		slog.Error("failed to get TV season", "error", err)
		common.WriteErrorResponse(c, "Failed to fetch TV season", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(c, result)
}

// handles TVDB trending movies route
func TrendingMoviesHandler(c *gin.Context) {
	slog.Debug("handling TVDB trending movies route")

	result, err := GetTrendingMovies(c.Request.Context())
	if err != nil {
		slog.Error("failed to get trending movies", "error", err)
		common.WriteErrorResponse(c, "Failed to fetch trending movies", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(c, result)
}

// handles TVDB search movies route
func SearchMoviesHandler(c *gin.Context) {
	slog.Debug("handling TVDB search movies route")

	query := c.Query("query")
	if query == "" {
		common.WriteErrorResponse(c, "query parameter is required", http.StatusBadRequest)
		return
	}

	result, err := SearchMovies(c.Request.Context(), query)
	if err != nil {
		slog.Error("failed to search movies", "error", err)
		common.WriteErrorResponse(c, "Failed to search movies", http.StatusInternalServerError)
		return
	}

	common.WriteJSONResponse(c, result)
}

// handles TVDB get movie route
func GetMovieHandler(c *gin.Context) {
	slog.Debug("handling TVDB get movie route")

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
