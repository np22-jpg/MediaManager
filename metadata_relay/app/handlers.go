package app

import (
	"log/slog"
	"net/http"
)

// RootHandler handles the root route and returns API information.
func RootHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling hello route", "method", r.Method, "path", r.URL.Path)
	JSONResponse(w, http.StatusOK, map[string]string{"message": "metadata_relay API. See docs for details."})
}
