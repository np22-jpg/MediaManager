package app

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// handles the root route
func RootHandler(c *gin.Context) {
	slog.Debug("handling hello route", "method", c.Request.Method, "path", c.Request.URL.Path)
	c.JSON(http.StatusOK, gin.H{"message": "metadata_relay API. See docs for details."})
}
