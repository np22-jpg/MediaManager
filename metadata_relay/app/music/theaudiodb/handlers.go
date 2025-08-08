package theaudiodb

import (
	"log/slog"
	"net/http"
	"strings"

	"relay/app/common"
)

var client *Client

// SetClient configures the internal client used by handlers.
func SetClient(c *Client) { client = c }

// GET /theaudiodb/artist?name=...
func SearchArtistHandler(w http.ResponseWriter, r *http.Request) {
	if client == nil {
		common.WriteErrorResponse(w, "TheAudioDB not configured", http.StatusServiceUnavailable)
		return
	}
	name := r.URL.Query().Get("name")
	if name == "" {
		common.WriteErrorResponse(w, "name parameter is required", http.StatusBadRequest)
		return
	}
	res, err := SearchArtist(r.Context(), name)
	if err != nil {
		slog.Error("theaudiodb search failed", "error", err)
		common.WriteErrorResponse(w, "upstream error", http.StatusBadGateway)
		return
	}
	common.WriteJSONResponse(w, res)
}

// GetURLParam extracts a URL parameter from the request path
func GetURLParam(r *http.Request, key string) string {
	// Get path parameters from context if they exist
	if params := r.Context().Value("pathParams"); params != nil {
		if paramMap, ok := params.(map[string]string); ok {
			return paramMap[key]
		}
	}

	// Fallback: extract from URL path - simple approach for IDs at the end
	path := r.URL.Path
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

// GET /theaudiodb/artist/{mbid} (MBID lookup with high-priority caching)
func GetArtistByMBIDHandler(w http.ResponseWriter, r *http.Request) {
	if client == nil {
		common.WriteErrorResponse(w, "TheAudioDB not configured", http.StatusServiceUnavailable)
		return
	}
	mbid := GetURLParam(r, "mbid")
	if mbid == "" {
		common.WriteErrorResponse(w, "mbid parameter is required", http.StatusBadRequest)
		return
	}
	res, err := GetArtistByMBID(r.Context(), mbid)
	if err != nil {
		slog.Error("theaudiodb artist mbid lookup failed", "error", err, "mbid", mbid)
		common.WriteErrorResponse(w, "upstream error", http.StatusBadGateway)
		return
	}
	common.WriteJSONResponse(w, res)
}

// GET /theaudiodb/album/{mbid} (MBID lookup with high-priority caching)
func GetAlbumByMBIDHandler(w http.ResponseWriter, r *http.Request) {
	if client == nil {
		common.WriteErrorResponse(w, "TheAudioDB not configured", http.StatusServiceUnavailable)
		return
	}
	mbid := GetURLParam(r, "mbid")
	if mbid == "" {
		common.WriteErrorResponse(w, "mbid parameter is required", http.StatusBadRequest)
		return
	}
	res, err := GetAlbumByMBID(r.Context(), mbid)
	if err != nil {
		slog.Error("theaudiodb album mbid lookup failed", "error", err, "mbid", mbid)
		common.WriteErrorResponse(w, "upstream error", http.StatusBadGateway)
		return
	}
	common.WriteJSONResponse(w, res)
}

// GET /theaudiodb/track/{mbid} (MBID lookup with high-priority caching)
func GetTrackByMBIDHandler(w http.ResponseWriter, r *http.Request) {
	if client == nil {
		common.WriteErrorResponse(w, "TheAudioDB not configured", http.StatusServiceUnavailable)
		return
	}
	mbid := GetURLParam(r, "mbid")
	if mbid == "" {
		common.WriteErrorResponse(w, "mbid parameter is required", http.StatusBadRequest)
		return
	}
	res, err := GetTrackByMBID(r.Context(), mbid)
	if err != nil {
		slog.Error("theaudiodb track mbid lookup failed", "error", err, "mbid", mbid)
		common.WriteErrorResponse(w, "upstream error", http.StatusBadGateway)
		return
	}
	common.WriteJSONResponse(w, res)
}
