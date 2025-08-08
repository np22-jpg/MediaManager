package app

import (
	"relay/app/anidb"
	"relay/app/jikan"
	"relay/app/music/musicbrainz"
	"relay/app/music/theaudiodb"
	"relay/app/seadex"
	"relay/app/tmdb"
	"relay/app/tvdb"
)

// RegisterRoutes registers all API routes
func RegisterRoutes(router *Router, musicBrainzEnabled bool, seadexEnabled bool, anidbEnabled bool, jikanEnabled bool) {
	// Root endpoint
	router.GET("/", RootHandler)

	// TMDB endpoints
	router.GET("/tmdb/tv/trending", tmdb.TrendingTVHandler)
	router.GET("/tmdb/tv/search", tmdb.SearchTVHandler)
	router.GET("/tmdb/tv/shows/{showId}", tmdb.GetTVShowHandler)
	router.GET("/tmdb/tv/shows/{showId}/{seasonNumber}", tmdb.GetTVSeasonHandler)
	router.GET("/tmdb/movies/trending", tmdb.TrendingMoviesHandler)
	router.GET("/tmdb/movies/search", tmdb.SearchMoviesHandler)
	router.GET("/tmdb/movies/{movieId}", tmdb.GetMovieHandler)

	// TVDB endpoints
	router.GET("/tvdb/tv/trending", tvdb.TrendingTVHandler)
	router.GET("/tvdb/tv/search", tvdb.SearchTVHandler)
	router.GET("/tvdb/tv/shows/{showId}", tvdb.GetTVShowHandler)
	router.GET("/tvdb/tv/seasons/{seasonId}", tvdb.GetTVSeasonHandler)
	router.GET("/tvdb/movies/trending", tvdb.TrendingMoviesHandler)
	router.GET("/tvdb/movies/search", tvdb.SearchMoviesHandler)
	router.GET("/tvdb/movies/{movieId}", tvdb.GetMovieHandler)

	// SeaDx endpoints (conditional)
	if seadexEnabled {
		router.GET("/seadx/search", seadex.SearchEntriesHandler)
		router.GET("/seadx/entries/{id}", seadex.GetEntryByIDHandler)
		router.GET("/seadx/anilist/{anilistId}", seadex.GetEntryByAnilistIDHandler)
		router.GET("/seadx/trending", seadex.GetTrendingEntriesHandler)
		router.GET("/seadx/release-groups", seadex.GetEntriesByReleaseGroupHandler)
		router.GET("/seadx/trackers", seadex.GetEntriesByTrackerHandler)
	}

	// AniDB endpoints (conditional)
	if anidbEnabled {
		router.GET("/anidb/anime/{id}", anidb.GetAnimeByIDHandler)
		router.GET("/anidb/hot", anidb.GetHotAnimeHandler)
		router.GET("/anidb/recommendations", anidb.GetRandomRecommendationHandler)
		router.GET("/anidb/similar", anidb.GetRandomSimilarHandler)
		router.GET("/anidb/main", anidb.GetMainPageDataHandler)
	}

	// Jikan endpoints (MyAnimeList API alternative - conditional)
	if jikanEnabled {
		router.GET("/jikan/anime/{id}", jikan.GetAnimeByIDHandler)
		router.GET("/jikan/top", jikan.GetTopAnimeHandler)
		router.GET("/jikan/seasonal", jikan.GetSeasonalAnimeHandler)
		router.GET("/jikan/search", jikan.SearchAnimeHandler)
		router.GET("/jikan/anime/{id}/recommendations", jikan.GetAnimeRecommendationsHandler)
		router.GET("/jikan/random", jikan.GetRandomAnimeHandler)
	}

	// TheAudioDB endpoints (independent of MusicBrainz)
	theaudiodb.RegisterRoutes(router)

	// MusicBrainz endpoints (conditional)
	if musicBrainzEnabled {
		router.GET("/musicbrainz/artists/search", musicbrainz.SearchArtistsHandler)
		router.GET("/musicbrainz/artists/search/advanced", musicbrainz.AdvancedSearchArtistsHandler)
		router.GET("/musicbrainz/artists/{mbid}", musicbrainz.GetArtistHandler)
		router.GET("/musicbrainz/artists/{mbid}/release-groups", musicbrainz.BrowseArtistReleaseGroupsHandler)
		router.GET("/musicbrainz/release-groups/search", musicbrainz.SearchReleaseGroupsHandler)
		router.GET("/musicbrainz/release-groups/{mbid}", musicbrainz.GetReleaseGroupHandler)
		router.GET("/musicbrainz/release-groups/{mbid}/releases", musicbrainz.BrowseReleaseGroupReleasesHandler)
		router.GET("/musicbrainz/releases/search", musicbrainz.SearchReleasesHandler)
		router.GET("/musicbrainz/releases/{mbid}", musicbrainz.GetReleaseHandler)
		router.GET("/musicbrainz/recordings/search", musicbrainz.SearchRecordingsHandler)
		router.GET("/musicbrainz/recordings/{mbid}", musicbrainz.GetRecordingHandler)
	}
}
