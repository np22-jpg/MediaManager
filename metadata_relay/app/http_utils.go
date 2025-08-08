package app

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"relay/app/metrics"
)

// HTTPError represents an HTTP error response
type HTTPError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const pathParamsKey contextKey = "pathParams"

// ResponseWriter wraps http.ResponseWriter with utility methods
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// StatusCode returns the captured status code
func (rw *ResponseWriter) StatusCode() int {
	if rw.statusCode == 0 {
		return http.StatusOK
	}
	return rw.statusCode
}

// JSONResponse sends a JSON response with the given status code and data
func JSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}

// ErrorResponse sends a JSON error response
func ErrorResponse(w http.ResponseWriter, status int, message string, err error) {
	response := HTTPError{
		Status:  status,
		Message: message,
	}
	if err != nil {
		response.Error = err.Error()
	}
	JSONResponse(w, status, response)
}

// GetURLParam extracts a URL parameter from the request path
// For paths like "/users/{id}", it extracts the id parameter
// This is a simple implementation that works with our current URL patterns
func GetURLParam(r *http.Request, key string) string {
	// Get path parameters from context if they exist
	if params := r.Context().Value("pathParams"); params != nil {
		if paramMap, ok := params.(map[string]string); ok {
			return paramMap[key]
		}
	}

	// Fallback: extract from URL path
	// This is a simplified approach - assumes parameter is at the end
	path := strings.TrimPrefix(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

// GetQueryParam gets a query parameter from the request
func GetQueryParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

// GetQueryParamInt gets a query parameter as an integer
func GetQueryParamInt(r *http.Request, key string, defaultValue int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}
	if parsed, err := strconv.Atoi(value); err == nil {
		return parsed
	}
	return defaultValue
}

// MiddlewareFunc represents HTTP middleware
type MiddlewareFunc func(http.Handler) http.Handler

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap the ResponseWriter to capture status code
			rw := &ResponseWriter{ResponseWriter: w}

			// Process request
			next.ServeHTTP(rw, r)

			// Log request details
			duration := time.Since(start)
			status := strconv.Itoa(rw.StatusCode())

			slog.Debug("HTTP request completed",
				"method", r.Method,
				"path", r.URL.Path,
				"status", status,
				"duration", duration,
			)
		})
	}
}

// MetricsMiddleware records HTTP metrics
func MetricsMiddleware() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap the ResponseWriter to capture status code
			rw := &ResponseWriter{ResponseWriter: w}

			// Process request
			next.ServeHTTP(rw, r)

			// Record metrics
			duration := time.Since(start)
			status := strconv.Itoa(rw.StatusCode())
			metrics.RecordHTTPRequest(r.Method, r.URL.Path, status, duration)
		})
	}
}

// RecoveryMiddleware recovers from panics and returns a 500 error
func RecoveryMiddleware() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					slog.Error("panic recovered", "error", err, "path", r.URL.Path, "method", r.Method)
					ErrorResponse(w, http.StatusInternalServerError, "Internal server error", nil)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// ChainMiddleware chains multiple middleware functions
func ChainMiddleware(middlewares ...MiddlewareFunc) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

// Router represents a simple HTTP router with path parameter support
type Router struct {
	mux        *http.ServeMux
	middleware []MiddlewareFunc
}

// NewRouter creates a new router
func NewRouter() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

// Use adds middleware to the router
func (r *Router) Use(middleware ...MiddlewareFunc) {
	r.middleware = append(r.middleware, middleware...)
}

// Handle registers a handler for the given pattern
func (r *Router) Handle(pattern string, handler http.Handler) {
	r.mux.Handle(pattern, handler)
}

// HandleFunc registers a handler function for the given pattern
func (r *Router) HandleFunc(pattern string, handler http.HandlerFunc) {
	// For patterns with parameters, we need special handling
	if strings.Contains(pattern, "{") {
		// Create a wrapper that extracts path parameters
		wrappedHandler := func(w http.ResponseWriter, req *http.Request) {
			params := r.extractParams(pattern, req.URL.Path)
			if len(params) > 0 {
				req = WithPathParams(req, params)
			}
			handler(w, req)
		}

		// Convert pattern to a format that http.ServeMux can handle
		wildcardPattern := r.convertToWildcard(pattern)
		r.mux.HandleFunc(wildcardPattern, wrappedHandler)
	} else {
		r.mux.HandleFunc(pattern, handler)
	}
}

// convertToWildcard converts {param} patterns to wildcard patterns
func (r *Router) convertToWildcard(pattern string) string {
	// For Go 1.22+ ServeMux, we can use wildcard patterns like "/users/{id}"
	// For older versions, we need to create catch-all patterns

	// Remove method prefix if present (e.g., "GET /users/{id}" -> "/users/{id}")
	if idx := strings.Index(pattern, " "); idx != -1 {
		pattern = pattern[idx+1:]
	}

	// Convert {param} to generic wildcard for simpler matching
	// This is a simplified approach - in production you might want gorilla/mux
	parts := strings.Split(pattern, "/")
	var result []string
	hasParam := false

	for _, part := range parts {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			hasParam = true
			break
		}
		result = append(result, part)
	}

	if hasParam && len(result) > 0 {
		// Create a catch-all pattern like "/users/"
		return strings.Join(result, "/") + "/"
	}

	return pattern
}

// extractParams extracts parameters from URL path based on pattern
func (r *Router) extractParams(pattern, path string) map[string]string {
	params := make(map[string]string)

	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	if len(patternParts) != len(pathParts) {
		return params
	}

	for i, part := range patternParts {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			key := strings.Trim(part, "{}")
			if i < len(pathParts) {
				params[key] = pathParts[i]
			}
		}
	}

	return params
}

// GET registers a GET handler
func (r *Router) GET(pattern string, handler http.HandlerFunc) {
	r.HandleFunc("GET "+pattern, handler)
}

// POST registers a POST handler
func (r *Router) POST(pattern string, handler http.HandlerFunc) {
	r.HandleFunc("POST "+pattern, handler)
}

// PUT registers a PUT handler
func (r *Router) PUT(pattern string, handler http.HandlerFunc) {
	r.HandleFunc("PUT "+pattern, handler)
}

// DELETE registers a DELETE handler
func (r *Router) DELETE(pattern string, handler http.HandlerFunc) {
	r.HandleFunc("DELETE "+pattern, handler)
}

// ServeHTTP implements the http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler := http.Handler(r.mux)
	for i := len(r.middleware) - 1; i >= 0; i-- {
		handler = r.middleware[i](handler)
	}
	handler.ServeHTTP(w, req)
}

// PathExtractor helps extract path parameters from URLs
// This is a simplified version - for production use, consider a more robust router
type PathExtractor struct {
	pattern string
}

// NewPathExtractor creates a new path extractor for the given pattern
func NewPathExtractor(pattern string) *PathExtractor {
	return &PathExtractor{pattern: pattern}
}

// Extract extracts parameters from the given path based on the pattern
func (pe *PathExtractor) Extract(path string) map[string]string {
	params := make(map[string]string)

	patternParts := strings.Split(strings.Trim(pe.pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	if len(patternParts) != len(pathParts) {
		return params
	}

	for i, part := range patternParts {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			key := strings.Trim(part, "{}")
			if i < len(pathParts) {
				params[key] = pathParts[i]
			}
		}
	}

	return params
}

// WithPathParams adds path parameters to the request context
func WithPathParams(r *http.Request, params map[string]string) *http.Request {
	ctx := context.WithValue(r.Context(), pathParamsKey, params)
	return r.WithContext(ctx)
}

// GetPathParam gets a path parameter from the request context
func GetPathParam(r *http.Request, key string) string {
	if params := r.Context().Value(pathParamsKey); params != nil {
		if paramMap, ok := params.(map[string]string); ok {
			return paramMap[key]
		}
	}
	return ""
}

// Static serves static files from the given directory
func Static(dir string) http.Handler {
	return http.StripPrefix("/", http.FileServer(http.Dir(dir)))
}
