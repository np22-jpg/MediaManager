package theaudiodb

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes wires TheAudioDB endpoints into the main router.
func RegisterRoutes(router *gin.Engine) {
	group := router.Group("/theaudiodb")
	{
		// Name-based search (existing)
		group.GET("/artist", SearchArtistHandler)

		// MBID-based lookups (new, with prioritized caching)
		group.GET("/artist/:mbid", GetArtistByMBIDHandler)
		group.GET("/album/:mbid", GetAlbumByMBIDHandler)
		group.GET("/track/:mbid", GetTrackByMBIDHandler)
	}
}
