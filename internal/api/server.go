package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"explorer451/internal/auth"
	"explorer451/internal/core"
	"explorer451/internal/models"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// CustomValidator is a custom validator for Echo
type CustomValidator struct {
	validator *validator.Validate
}

// Validate validates the struct
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// Server represents the HTTP server
type Server struct {
	echo *echo.Echo
	core *core.Core
	auth *auth.Auth
}

// NewServer creates a new HTTP server
func NewServer(core *core.Core) (*Server, error) {
	s := &Server{
		echo: echo.New(),
		core: core,
	}

	// Initialize authentication if database is available
	if core.DB != nil {
		callbacks := &auth.Callbacks{
			GetUser: func(id string) (*models.User, error) {
				ctx := context.Background()
				return core.Repository.GetUser(ctx, id)
			},
		}

		authSystem, err := auth.New(context.Background(), core.Repository, *core.Config, core.DB.DB, callbacks, core.Logger)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize authentication: %w", err)
		}
		s.auth = authSystem
	}

	// Configure middleware
	s.echo.Use(middleware.Recover())
	s.echo.Use(middleware.Logger())
	s.echo.Use(middleware.CORS())
	s.echo.Use(middleware.RequestID())
	s.echo.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 30 * time.Second,
	}))

	// Setup validator
	s.echo.Validator = &CustomValidator{validator: validator.New()}

	// Setup routes
	s.setupRoutes()

	return s, nil
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

	// Authentication routes (if auth is available)
	if s.auth != nil {
		auth := api.Group("/auth")
		auth.POST("/login", s.auth.Login)
		auth.POST("/logout", s.auth.LogoutHandler)
		auth.GET("/profile", s.auth.Profile, s.auth.Middleware)
		auth.GET("/oidc/login", s.auth.OIDCLogin)
		auth.GET("/oidc/callback", s.auth.OIDCCallback)

		// Protected API endpoints (require authentication)
		protected := api.Group("", s.auth.Middleware)
		protected.GET("/buckets", s.listBuckets)
		protected.GET("/buckets/:bucket/details", s.getBucketDetails)
		protected.GET("/buckets/:bucket/objects", s.listObjects)
		protected.GET("/buckets/:bucket/objects/*", s.getPresignedURL)
		protected.HEAD("/buckets/:bucket/objects/*", s.getObjectMetadata)
		protected.DELETE("/buckets/:bucket/objects/*", s.deleteObject)
		protected.POST("/buckets/:bucket/objects", s.createFolder)
		protected.POST("/buckets/:bucket/presigned-post-url", s.generatePresignedPostURL)
	} else {
		// Fallback: unprotected routes if no auth system
		api.GET("/buckets", s.listBuckets)
		api.GET("/buckets/:bucket/details", s.getBucketDetails)
		api.GET("/buckets/:bucket/objects", s.listObjects)
		api.GET("/buckets/:bucket/objects/*", s.getPresignedURL)
		api.HEAD("/buckets/:bucket/objects/*", s.getObjectMetadata)
		api.DELETE("/buckets/:bucket/objects/*", s.deleteObject)
		api.POST("/buckets/:bucket/objects", s.createFolder)
		api.POST("/buckets/:bucket/presigned-post-url", s.generatePresignedPostURL)
	}
}
