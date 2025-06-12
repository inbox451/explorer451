package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"explorer451/internal/models"

	"github.com/aws/smithy-go"
	"github.com/labstack/echo/v4"
)

// listBuckets handles GET /api/buckets
func (s *Server) listBuckets(c echo.Context) error {
	buckets, err := s.core.S3Service.ListBuckets(c.Request().Context())
	if err != nil {
		s.core.Logger.Error().Err(err).Msg("Error listing buckets")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list buckets")
	}

	return c.JSON(http.StatusOK, buckets)
}

// listObjects handles GET /api/buckets/:bucket/objects
func (s *Server) listObjects(c echo.Context) error {
	bucket := c.Param("bucket")
	prefix := c.QueryParam("prefix")
	nextToken := c.QueryParam("nextToken")
	delimiter := c.QueryParam("delimiter")

	// Parse maxKeys parameter if provided
	maxKeys := int32(1000) // Default
	if c.QueryParam("maxKeys") != "" {
		if val, err := strconv.ParseInt(c.QueryParam("maxKeys"), 10, 32); err == nil {
			maxKeys = int32(val)
		}
	}

	objects, err := s.core.S3Service.ListObjects(
		c.Request().Context(),
		bucket,
		prefix,
		nextToken,
		delimiter,
		maxKeys,
	)

	if err != nil {
		// Map common AWS errors to appropriate HTTP status
		if isNoSuchBucketError(err) {
			return echo.NewHTTPError(http.StatusNotFound, "Bucket not found")
		}
		if isAccessDeniedError(err) {
			return echo.NewHTTPError(http.StatusForbidden, "Access denied")
		}

		s.core.Logger.Error().Err(err).Str("bucket", bucket).Msg("Error listing objects")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list objects")
	}

	return c.JSON(http.StatusOK, objects)
}

// getPresignedURL handles GET /api/buckets/:bucket/objects/:key/presigned-url
func (s *Server) getPresignedURL(c echo.Context) error {
	bucket := c.Param("bucket")
	key := c.Param("key")

	// Parse expiration time in seconds (default 15 minutes if not specified)
	expiresIn := int64(15 * 60)
	if c.QueryParam("expiresIn") != "" {
		if val, err := strconv.ParseInt(c.QueryParam("expiresIn"), 10, 64); err == nil {
			expiresIn = val
		}
	}

	url, err := s.core.S3Service.GetPresignedURL(c.Request().Context(), bucket, key, expiresIn)
	if err != nil {
		if isNoSuchBucketError(err) {
			return echo.NewHTTPError(http.StatusNotFound, "Bucket not found")
		}
		if isNoSuchKeyError(err) {
			return echo.NewHTTPError(http.StatusNotFound, "Object not found")
		}
		if isAccessDeniedError(err) {
			return echo.NewHTTPError(http.StatusForbidden, "Access denied")
		}

		s.core.Logger.Error().
			Err(err).
			Str("bucket", bucket).
			Str("key", key).
			Msg("Error generating presigned URL")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate presigned URL")
	}

	return c.JSON(http.StatusOK, map[string]string{"url": url})
}

// deleteObject handles DELETE /api/buckets/:bucket/objects/*
func (s *Server) deleteObject(c echo.Context) error {
	bucket := c.Param("bucket")
	key := c.Param("*")
	recursive := c.QueryParam("recursive") == "true"

	// If recursive is true, delete by prefix (folder deletion)
	if recursive {
		err := s.core.S3Service.DeleteObjectsByPrefix(c.Request().Context(), bucket, key)
		if err != nil {
			if isNoSuchBucketError(err) {
				return echo.NewHTTPError(http.StatusNotFound, "Bucket not found")
			}
			if isAccessDeniedError(err) {
				return echo.NewHTTPError(http.StatusForbidden, "Access denied")
			}

			s.core.Logger.Error().
				Err(err).
				Str("bucket", bucket).
				Str("prefix", key).
				Msg("Error deleting objects by prefix")
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete folder")
		}
	} else {
		// Single object deletion
		err := s.core.S3Service.DeleteObject(c.Request().Context(), bucket, key)
		if err != nil {
			if isNoSuchBucketError(err) {
				return echo.NewHTTPError(http.StatusNotFound, "Bucket not found")
			}
			if isNoSuchKeyError(err) {
				return echo.NewHTTPError(http.StatusNotFound, "Object not found")
			}
			if isAccessDeniedError(err) {
				return echo.NewHTTPError(http.StatusForbidden, "Access denied")
			}

			s.core.Logger.Error().
				Err(err).
				Str("bucket", bucket).
				Str("key", key).
				Msg("Error deleting object")
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete object")
		}
	}

	return c.NoContent(http.StatusNoContent)
}

// createFolder handles POST /api/buckets/:bucket/objects
func (s *Server) createFolder(c echo.Context) error {
	bucket := c.Param("bucket")

	var req models.CreateFolderRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Validate that key ends with '/' or add it if missing
	if !strings.HasSuffix(req.Key, "/") {
		req.Key = req.Key + "/"
	}

	// Validate type field
	if req.Type != "folder" {
		return echo.NewHTTPError(http.StatusBadRequest, "Type must be 'folder'")
	}

	err := s.core.S3Service.CreateFolder(c.Request().Context(), bucket, req.Key)
	if err != nil {
		if isNoSuchBucketError(err) {
			return echo.NewHTTPError(http.StatusNotFound, "Bucket not found")
		}
		if isAccessDeniedError(err) {
			return echo.NewHTTPError(http.StatusForbidden, "Access denied")
		}

		s.core.Logger.Error().
			Err(err).
			Str("bucket", bucket).
			Str("key", req.Key).
			Msg("Error creating folder")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create folder")
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "Folder created successfully",
		"key":     req.Key,
	})
}

// generatePresignedPostURL handles POST /api/buckets/:bucket/presigned-post-url
func (s *Server) generatePresignedPostURL(c echo.Context) error {
	bucket := c.Param("bucket")

	var req models.PresignedPostURLRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Validate required fields
	if req.Key == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Key is required")
	}
	if req.ContentType == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Content type is required")
	}

	// Set default values if not provided
	expiresIn := time.Duration(req.ExpiresInSeconds) * time.Second
	if req.ExpiresInSeconds <= 0 {
		expiresIn = 15 * time.Minute // Default to 15 minutes
	}

	maxSize := req.MaxSizeBytes
	if maxSize <= 0 {
		maxSize = 10 * 1024 * 1024 // Default to 10MB
	}

	response, err := s.core.S3Service.GeneratePresignedPostURL(
		c.Request().Context(),
		bucket,
		req.Key,
		req.ContentType,
		expiresIn,
		maxSize,
	)
	if err != nil {
		if isNoSuchBucketError(err) {
			return echo.NewHTTPError(http.StatusNotFound, "Bucket not found")
		}
		if isAccessDeniedError(err) {
			return echo.NewHTTPError(http.StatusForbidden, "Access denied")
		}

		s.core.Logger.Error().
			Err(err).
			Str("bucket", bucket).
			Str("key", req.Key).
			Msg("Error generating presigned POST URL")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate presigned POST URL")
	}

	return c.JSON(http.StatusOK, response)
}

// Helper functions to identify AWS error types
func isNoSuchBucketError(err error) bool {
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		return apiErr.ErrorCode() == "NoSuchBucket"
	}
	return false
}

func isAccessDeniedError(err error) bool {
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		return apiErr.ErrorCode() == "AccessDenied"
	}
	return false
}

func isNoSuchKeyError(err error) bool {
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		return apiErr.ErrorCode() == "NoSuchKey"
	}
	return false
}
