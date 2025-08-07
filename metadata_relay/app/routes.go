package app

import (
	"relay/app/tmdb"
	"relay/app/tvdb"

	"github.com/gin-gonic/gin"
)

// registers all API routes
func RegisterRoutes(router *gin.Engine) {
	// Root endpoint
	router.GET("/", RootHandler)

	// TMDB endpoints group
	tmdbGroup := router.Group("/tmdb")
	{
		// TV endpoints
		tvGroup := tmdbGroup.Group("/tv")
		{
			tvGroup.GET("/trending", tmdb.TrendingTVHandler)
			tvGroup.GET("/search", tmdb.SearchTVHandler)
			tvGroup.GET("/shows/:showId", tmdb.GetTVShowHandler)
			tvGroup.GET("/shows/:showId/:seasonNumber", tmdb.GetTVSeasonHandler)
		}

		// Movie endpoints
		moviesGroup := tmdbGroup.Group("/movies")
		{
			moviesGroup.GET("/trending", tmdb.TrendingMoviesHandler)
			moviesGroup.GET("/search", tmdb.SearchMoviesHandler)
			moviesGroup.GET("/:movieId", tmdb.GetMovieHandler)
		}
	}

	// TVDB endpoints group
	tvdbGroup := router.Group("/tvdb")
	{
		// TV endpoints
		tvGroup := tvdbGroup.Group("/tv")
		{
			tvGroup.GET("/trending", tvdb.TrendingTVHandler)
			tvGroup.GET("/search", tvdb.SearchTVHandler)
			tvGroup.GET("/shows/:showId", tvdb.GetTVShowHandler)
		}

		// Season endpoints
		tvGroup.GET("/seasons/:seasonId", tvdb.GetTVSeasonHandler)

		// Movie endpoints
		moviesGroup := tvdbGroup.Group("/movies")
		{
			moviesGroup.GET("/trending", tvdb.TrendingMoviesHandler)
			moviesGroup.GET("/search", tvdb.SearchMoviesHandler)
			moviesGroup.GET("/:movieId", tvdb.GetMovieHandler)
		}
	}
}
