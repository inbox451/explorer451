package api

import (
	"errors"
	"net/http"
	"strconv"

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
