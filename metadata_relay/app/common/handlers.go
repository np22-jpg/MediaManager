package common

import (
	"github.com/gin-gonic/gin"
)

// WriteJSONResponse writes a JSON response with HTTP 200 status.
func WriteJSONResponse(c *gin.Context, data any) {
	c.JSON(200, data)
}

// WriteErrorResponse writes an error response with the specified status code.
func WriteErrorResponse(c *gin.Context, message string, statusCode int) {
	c.JSON(statusCode, gin.H{"error": message})
}
