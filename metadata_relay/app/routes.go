package app

import "net/http"

func RegisterRoutes(mux *http.ServeMux) {
	// Root endpoint
	mux.HandleFunc("/", HelloHandler)

	// TMDB endpoints
	mux.HandleFunc("GET /tmdb/tv/trending", TMDBTrendingTVHandler)
	mux.HandleFunc("GET /tmdb/tv/search", TMDBSearchTVHandler)
	mux.HandleFunc("GET /tmdb/tv/shows/{showId}", TMDBGetTVShowHandler)
	mux.HandleFunc("GET /tmdb/tv/shows/{showId}/{seasonNumber}", TMDBGetTVSeasonHandler)
	mux.HandleFunc("GET /tmdb/movies/trending", TMDBTrendingMoviesHandler)
	mux.HandleFunc("GET /tmdb/movies/search", TMDBSearchMoviesHandler)
	mux.HandleFunc("GET /tmdb/movies/{movieId}", TMDBGetMovieHandler)

	// TVDB endpoints
	mux.HandleFunc("GET /tvdb/tv/trending", TVDBTrendingTVHandler)
	mux.HandleFunc("GET /tvdb/tv/search", TVDBSearchTVHandler)
	mux.HandleFunc("GET /tvdb/tv/shows/{showId}", TVDBGetTVShowHandler)
	mux.HandleFunc("GET /tvdb/tv/seasons/{seasonId}", TVDBGetTVSeasonHandler)
	mux.HandleFunc("GET /tvdb/movies/trending", TVDBTrendingMoviesHandler)
	mux.HandleFunc("GET /tvdb/movies/search", TVDBSearchMoviesHandler)
	mux.HandleFunc("GET /tvdb/movies/{movieId}", TVDBGetMovieHandler)
}
