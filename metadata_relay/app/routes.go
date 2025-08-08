package app

import (
	"relay/app/anidb"
	"relay/app/music/musicbrainz"
	"relay/app/music/theaudiodb"
	"relay/app/seadex"
	"relay/app/tmdb"
	"relay/app/tvdb"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all API routes
func RegisterRoutes(router *gin.Engine, musicBrainzEnabled bool, seadexEnabled bool, anidbEnabled bool) {
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

	// SeaDx endpoints group (conditional)
	if seadexEnabled {
		seadexGroup := router.Group("/seadx")
		{
			seadexGroup.GET("/search", seadex.SearchEntriesHandler)
			seadexGroup.GET("/entries/:id", seadex.GetEntryByIDHandler)
			seadexGroup.GET("/anilist/:anilistId", seadex.GetEntryByAnilistIDHandler)
			seadexGroup.GET("/trending", seadex.GetTrendingEntriesHandler)
			seadexGroup.GET("/release-groups", seadex.GetEntriesByReleaseGroupHandler)
			seadexGroup.GET("/trackers", seadex.GetEntriesByTrackerHandler)
		}
	}

	// AniDB endpoints group (conditional)
	if anidbEnabled {
		anidbGroup := router.Group("/anidb")
		{
			anidbGroup.GET("/anime/:id", anidb.GetAnimeByIDHandler)
			anidbGroup.GET("/hot", anidb.GetHotAnimeHandler)
			anidbGroup.GET("/recommendations", anidb.GetRandomRecommendationHandler)
			anidbGroup.GET("/similar", anidb.GetRandomSimilarHandler)
			anidbGroup.GET("/main", anidb.GetMainPageDataHandler)
		}
	}

	// TheAudioDB endpoints (independent of MusicBrainz)
	theaudiodb.RegisterRoutes(router)

	// MusicBrainz endpoints group (conditional)
	if musicBrainzEnabled {
		musicbrainzGroup := router.Group("/musicbrainz")
		{
			// Artist endpoints
			artistGroup := musicbrainzGroup.Group("/artists")
			{
				artistGroup.GET("/search", musicbrainz.SearchArtistsHandler)
				artistGroup.GET("/search/advanced", musicbrainz.AdvancedSearchArtistsHandler)
				artistGroup.GET("/:mbid", musicbrainz.GetArtistHandler)
				artistGroup.GET("/:mbid/release-groups", musicbrainz.BrowseArtistReleaseGroupsHandler)
			}

			// Release Group endpoints
			releaseGroupGroup := musicbrainzGroup.Group("/release-groups")
			{
				releaseGroupGroup.GET("/search", musicbrainz.SearchReleaseGroupsHandler)
				releaseGroupGroup.GET("/:mbid", musicbrainz.GetReleaseGroupHandler)
				releaseGroupGroup.GET("/:mbid/releases", musicbrainz.BrowseReleaseGroupReleasesHandler)
			}

			// Release endpoints
			releaseGroup := musicbrainzGroup.Group("/releases")
			{
				releaseGroup.GET("/search", musicbrainz.SearchReleasesHandler)
				releaseGroup.GET("/:mbid", musicbrainz.GetReleaseHandler)
			}

			// Recording endpoints
			recordingGroup := musicbrainzGroup.Group("/recordings")
			{
				recordingGroup.GET("/search", musicbrainz.SearchRecordingsHandler)
				recordingGroup.GET("/:mbid", musicbrainz.GetRecordingHandler)
			}
		}
	}
}
