package app

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine) {
	// Root endpoint
	router.GET("/", HelloHandler)

	// TMDB endpoints group
	tmdbGroup := router.Group("/tmdb")
	{
		// TV endpoints
		tvGroup := tmdbGroup.Group("/tv")
		{
			tvGroup.GET("/trending", TMDBTrendingTVHandler)
			tvGroup.GET("/search", TMDBSearchTVHandler)
			tvGroup.GET("/shows/:showId", TMDBGetTVShowHandler)
			tvGroup.GET("/shows/:showId/:seasonNumber", TMDBGetTVSeasonHandler)
		}

		// Movie endpoints
		moviesGroup := tmdbGroup.Group("/movies")
		{
			moviesGroup.GET("/trending", TMDBTrendingMoviesHandler)
			moviesGroup.GET("/search", TMDBSearchMoviesHandler)
			moviesGroup.GET("/:movieId", TMDBGetMovieHandler)
		}
	}

	// TVDB endpoints group
	tvdbGroup := router.Group("/tvdb")
	{
		// TV endpoints
		tvGroup := tvdbGroup.Group("/tv")
		{
			tvGroup.GET("/trending", TVDBTrendingTVHandler)
			tvGroup.GET("/search", TVDBSearchTVHandler)
			tvGroup.GET("/shows/:showId", TVDBGetTVShowHandler)
		}

		// Season endpoints
		tvGroup.GET("/seasons/:seasonId", TVDBGetTVSeasonHandler)

		// Movie endpoints
		moviesGroup := tvdbGroup.Group("/movies")
		{
			moviesGroup.GET("/trending", TVDBTrendingMoviesHandler)
			moviesGroup.GET("/search", TVDBSearchMoviesHandler)
			moviesGroup.GET("/:movieId", TVDBGetMovieHandler)
		}
	}
}
