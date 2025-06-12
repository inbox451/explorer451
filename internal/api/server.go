package api

import (
	"context"
	"net/http"
	"time"

	"explorer451/internal/core"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Server represents the HTTP server
type Server struct {
	echo *echo.Echo
	core *core.Core
}

// NewServer creates a new HTTP server
func NewServer(core *core.Core) *Server {
	s := &Server{
		echo: echo.New(),
		core: core,
	}

	// Configure middleware
	s.echo.Use(middleware.Recover())
	s.echo.Use(middleware.Logger())
	s.echo.Use(middleware.CORS())
	s.echo.Use(middleware.RequestID())
	s.echo.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 30 * time.Second,
	}))

	// Setup routes
	s.setupRoutes()

	return s
}

// Start starts the HTTP server
func (s *Server) Start(address string) error {
	return s.echo.Start(address)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}

// setupRoutes configures all API routes
func (s *Server) setupRoutes() {
	// Health check endpoint
	s.echo.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// API endpoints
	api := s.echo.Group("/api")

	// Bucket endpoints
	api.GET("/buckets", s.listBuckets)
	api.GET("/buckets/:bucket/objects", s.listObjects)
	api.GET("/buckets/:bucket/objects/*", s.getPresignedURL)
	api.DELETE("/buckets/:bucket/objects/*", s.deleteObject)
	api.POST("/buckets/:bucket/objects", s.createFolder)
	api.POST("/buckets/:bucket/presigned-post-url", s.generatePresignedPostURL)
}
