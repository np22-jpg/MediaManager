package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"relay/app"
	"relay/app/metrics"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// getLogLevel parses the VERBOSITY environment variable and returns the appropriate slog.Level
func getLogLevel() slog.Level {
	verbosity := strings.ToLower(os.Getenv("LOG_LEVEL"))
	switch verbosity {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo // Default to info level
	}
}

// httpMetricsMiddleware wraps handlers to record HTTP metrics
func httpMetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Create a response writer wrapper to capture status code
		ww := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}
		
		// Serve the request
		next.ServeHTTP(ww, r)
		
		// Record metrics
		duration := time.Since(start)
		status := strconv.Itoa(ww.statusCode)
		metrics.RecordHTTPRequest(r.Method, r.URL.Path, status, duration)
		
		slog.Debug("HTTP request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", ww.statusCode,
			"duration", duration,
			"remote_addr", r.RemoteAddr,
		)
	})
}

// responseWriterWrapper wraps http.ResponseWriter to capture status code
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func main() {
	// Set up structured logger with configurable level
	logLevel := getLogLevel()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	slog.Info("starting server")

	// Build the router
	mux := http.NewServeMux()
	app.RegisterRoutes(mux)

	// Add Prometheus metrics endpoint (without middleware to avoid self-monitoring)
	mux.Handle("/metrics", promhttp.Handler())

	// Wrap with metrics middleware
	handler := httpMetricsMiddleware(mux)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in goroutine
	go func() {
		slog.Info("listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Handle shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	slog.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
	} else {
		slog.Info("server stopped gracefully")
	}
}
