// main.go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type S3Handler struct {
	s3Client *s3.Client
}

type ListBucketResponse struct {
	Items                 []S3Item `json:"items"`
	NextContinuationToken string   `json:"nextContinuationToken,omitempty"`
	IsTruncated           bool     `json:"isTruncated"`
	TotalItems            int32    `json:"totalItems"`
	PageSize              int32    `json:"pageSize"`
}

type S3Item struct {
	Key          string `json:"key"`
	Size         int64  `json:"size"`
	LastModified string `json:"lastModified,omitempty"`
	IsFolder     bool   `json:"isFolder"`
	Type         string `json:"type"` // "folder" or "file"
	ContentType  string `json:"contentType,omitempty"`
}

func NewS3Handler() (*S3Handler, error) {
	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		awsRegion = "eu-central-1"
		fmt.Printf("Defaulting ro AWS Region %s\n", awsRegion)
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(awsRegion))
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg)
	return &S3Handler{s3Client: s3Client}, nil
}

func (h *S3Handler) ListBucketContents(c echo.Context) error {
	bucket := c.Param("bucket")
	prefix := c.QueryParam("prefix")
	continuationToken := c.QueryParam("continuationToken")
	delimiter := "/" // Using delimiter to get folder-like hierarchy

	// Parse page size from query parameter, default to 100
	pageSize := int32(100)
	if sizeParm := c.QueryParam("pageSize"); sizeParm != "" {
		if size, err := strconv.ParseInt(sizeParm, 10, 32); err == nil {
			pageSize = int32(size)
		}
	}

	input := &s3.ListObjectsV2Input{
		Bucket:            &bucket,
		Prefix:            &prefix,
		MaxKeys:           &pageSize,
		ContinuationToken: &continuationToken,
		Delimiter:         &delimiter, // This will group objects by "folder"
	}

	// If no continuation token is provided, don't include it in the request
	if continuationToken == "" {
		input.ContinuationToken = nil
	}

	result, err := h.s3Client.ListObjectsV2(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	items := make([]S3Item, 0)

	// Add common prefixes (folders)
	for _, prefix := range result.CommonPrefixes {
		if prefix.Prefix == nil {
			continue
		}
		folderName := *prefix.Prefix
		// Remove trailing delimiter if present
		if strings.HasSuffix(folderName, delimiter) {
			folderName = folderName[:len(folderName)-1]
		}

		items = append(items, S3Item{
			Key:      folderName,
			IsFolder: true,
			Type:     "folder",
			Size:     0, // Folders don't have a size
		})
	}

	// Add objects (files)
	for _, obj := range result.Contents {
		if obj.Key == nil {
			continue
		}
		// Skip objects that represent folders (ending with delimiter)
		if strings.HasSuffix(*obj.Key, delimiter) {
			continue
		}

		contentType := "application/octet-stream" // default content type
		// Get content type for the object
		headInput := &s3.HeadObjectInput{
			Bucket: &bucket,
			Key:    obj.Key,
		}
		if headOutput, err := h.s3Client.HeadObject(c.Request().Context(), headInput); err == nil && headOutput.ContentType != nil {
			contentType = *headOutput.ContentType
		}

		item := S3Item{
			Key:         *obj.Key,
			IsFolder:    false,
			Type:        "file",
			ContentType: contentType,
		}

		// Handle optional/pointer fields
		if obj.Size != nil {
			item.Size = *obj.Size
		}
		if obj.LastModified != nil {
			item.LastModified = obj.LastModified.Format("2006-01-02T15:04:05Z07:00")
		}

		items = append(items, item)
	}

	response := ListBucketResponse{
		Items:       items,
		IsTruncated: false, // Default value
		TotalItems:  0,     // Default value
		PageSize:    pageSize,
	}

	// Handle optional/pointer fields from result
	if result.IsTruncated != nil {
		response.IsTruncated = *result.IsTruncated
	}
	if result.KeyCount != nil {
		response.TotalItems = *result.KeyCount
	}
	if result.IsTruncated != nil && *result.IsTruncated && result.NextContinuationToken != nil {
		response.NextContinuationToken = *result.NextContinuationToken
	}

	return c.JSON(http.StatusOK, response)
}

func (h *S3Handler) ListBuckets(c echo.Context) error {
	result, err := h.s3Client.ListBuckets(c.Request().Context(), &s3.ListBucketsInput{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	buckets := make([]string, 0)
	for _, bucket := range result.Buckets {
		if bucket.Name != nil {
			buckets = append(buckets, *bucket.Name)
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"buckets": buckets,
	})
}

func main() {
	// Create new Echo instance
	e := echo.New()
	e.HideBanner = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Create S3 handler
	s3Handler, err := NewS3Handler()
	if err != nil {
		log.Fatalf("Failed to create S3 handler: %v", err)
	}

	// Routes
	e.GET("/api/buckets", s3Handler.ListBuckets)
	e.GET("/api/buckets/:bucket/objects", s3Handler.ListBucketContents)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
