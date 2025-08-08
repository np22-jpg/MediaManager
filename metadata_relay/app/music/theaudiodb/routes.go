package theaudiodb

import (
	"net/http"
)

// RouterInterface represents the interface for our custom router
type RouterInterface interface {
	GET(pattern string, handler http.HandlerFunc)
}

// RegisterRoutes wires TheAudioDB endpoints into the main router.
func RegisterRoutes(router RouterInterface) {
	// Name-based search (existing)
	router.GET("/theaudiodb/artist", SearchArtistHandler)

	// MBID-based lookups (new, with prioritized caching)
	router.GET("/theaudiodb/artist/{mbid}", GetArtistByMBIDHandler)
	router.GET("/theaudiodb/album/{mbid}", GetAlbumByMBIDHandler)
	router.GET("/theaudiodb/track/{mbid}", GetTrackByMBIDHandler)
}
