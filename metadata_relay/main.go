package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"relay/app"
	"relay/app/cache"
	"relay/app/metrics"
	"relay/app/tmdb"
	"relay/app/tvdb"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ginMetricsMiddleware wraps Gin handlers to record HTTP metrics
func ginMetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Record metrics
		duration := time.Since(start)
		status := strconv.Itoa(c.Writer.Status())
		metrics.RecordHTTPRequest(c.Request.Method, c.FullPath(), status, duration)

		slog.Debug("HTTP request completed",
			"method", c.Request.Method,
			"path", c.FullPath(),
			"status", c.Writer.Status(),
			"duration", duration,
			"remote_addr", c.ClientIP(),
		)
	}
}

func main() {
	// Load configuration
	if err := LoadConfig(); err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Initialize cache
	cache.InitCache(AppConfig.CacheHost, AppConfig.CachePort, AppConfig.CacheDB)

	// Initialize API clients
	tmdb.InitTMDB(AppConfig.TMDBAPIKey)
	tvdb.InitTVDB(AppConfig.TVDBAPIKey)

	// Set up structured logger with configurable level
	logLevel := AppConfig.GetLogLevel()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	slog.Info("starting server")

	// Set Gin mode based on log level
	if logLevel == slog.LevelDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	router := gin.New()

	// Add Gin's default middleware
	router.Use(gin.Logger(), gin.Recovery())

	// Add our metrics middleware
	router.Use(ginMetricsMiddleware())

	// Register routes
	app.RegisterRoutes(router)

	// Add Prometheus metrics endpoint (without middleware to avoid self-monitoring)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	srv := &http.Server{
		Addr:         AppConfig.GetServerAddr(),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in goroutine
	go func() {
		slog.Info("listening on " + AppConfig.GetServerAddr())
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
