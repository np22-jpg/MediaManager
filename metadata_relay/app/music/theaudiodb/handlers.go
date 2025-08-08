package theaudiodb

import (
	"log/slog"
	"net/http"

	"relay/app/common"

	"github.com/gin-gonic/gin"
)

var client *Client

// SetClient configures the internal client used by handlers.
func SetClient(c *Client) { client = c }

// GET /theaudiodb/artist?name=...
func SearchArtistHandler(c *gin.Context) {
	if client == nil {
		common.WriteErrorResponse(c, "TheAudioDB not configured", http.StatusServiceUnavailable)
		return
	}
	name := c.Query("name")
	if name == "" {
		common.WriteErrorResponse(c, "name parameter is required", http.StatusBadRequest)
		return
	}
	res, err := SearchArtist(c.Request.Context(), name)
	if err != nil {
		slog.Error("theaudiodb search failed", "error", err)
		common.WriteErrorResponse(c, "upstream error", http.StatusBadGateway)
		return
	}
	common.WriteJSONResponse(c, res)
}

// GET /theaudiodb/artist/:mbid (MBID lookup with high-priority caching)
func GetArtistByMBIDHandler(c *gin.Context) {
	if client == nil {
		common.WriteErrorResponse(c, "TheAudioDB not configured", http.StatusServiceUnavailable)
		return
	}
	mbid := c.Param("mbid")
	if mbid == "" {
		common.WriteErrorResponse(c, "mbid parameter is required", http.StatusBadRequest)
		return
	}
	res, err := GetArtistByMBID(c.Request.Context(), mbid)
	if err != nil {
		slog.Error("theaudiodb artist mbid lookup failed", "error", err, "mbid", mbid)
		common.WriteErrorResponse(c, "upstream error", http.StatusBadGateway)
		return
	}
	common.WriteJSONResponse(c, res)
}

// GET /theaudiodb/album/:mbid (MBID lookup with high-priority caching)
func GetAlbumByMBIDHandler(c *gin.Context) {
	if client == nil {
		common.WriteErrorResponse(c, "TheAudioDB not configured", http.StatusServiceUnavailable)
		return
	}
	mbid := c.Param("mbid")
	if mbid == "" {
		common.WriteErrorResponse(c, "mbid parameter is required", http.StatusBadRequest)
		return
	}
	res, err := GetAlbumByMBID(c.Request.Context(), mbid)
	if err != nil {
		slog.Error("theaudiodb album mbid lookup failed", "error", err, "mbid", mbid)
		common.WriteErrorResponse(c, "upstream error", http.StatusBadGateway)
		return
	}
	common.WriteJSONResponse(c, res)
}

// GET /theaudiodb/track/:mbid (MBID lookup with high-priority caching)
func GetTrackByMBIDHandler(c *gin.Context) {
	if client == nil {
		common.WriteErrorResponse(c, "TheAudioDB not configured", http.StatusServiceUnavailable)
		return
	}
	mbid := c.Param("mbid")
	if mbid == "" {
		common.WriteErrorResponse(c, "mbid parameter is required", http.StatusBadRequest)
		return
	}
	res, err := GetTrackByMBID(c.Request.Context(), mbid)
	if err != nil {
		slog.Error("theaudiodb track mbid lookup failed", "error", err, "mbid", mbid)
		common.WriteErrorResponse(c, "upstream error", http.StatusBadGateway)
		return
	}
	common.WriteJSONResponse(c, res)
}
