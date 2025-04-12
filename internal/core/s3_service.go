package core

import (
	"context"
	"explorer451/internal/models"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Service handles S3 operations
type S3Service struct {
	core *Core
}

// NewS3Service creates a new S3Service
func NewS3Service(core *Core) *S3Service {
	return &S3Service{
		core: core,
	}
}

// ListBuckets lists all S3 buckets the caller has access to
func (s *S3Service) ListBuckets(ctx context.Context) ([]models.Bucket, error) {
	s.core.Logger.Debug().Msg("Listing buckets")

	output, err := s.core.S3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		s.core.Logger.Error().Err(err).Msg("Failed to list buckets")
		return nil, err
	}

	buckets := make([]models.Bucket, len(output.Buckets))
	for i, b := range output.Buckets {
		buckets[i] = models.Bucket{
			Name:         aws.ToString(b.Name),
			CreationDate: aws.ToTime(b.CreationDate),
		}
	}

	return buckets, nil
}

// ListObjects lists objects in a bucket with optional prefix for folder navigation
func (s *S3Service) ListObjects(ctx context.Context, bucket, prefix, nextToken string, delimiter string, maxKeys int32) (*models.ListObjectsResponse, error) {
	s.core.Logger.Debug().
		Str("bucket", bucket).
		Str("prefix", prefix).
		Str("nextToken", nextToken).
		Msg("Listing objects")

	if delimiter == "" {
		delimiter = "/" // Default delimiter for folder-like navigation
	}

	if maxKeys <= 0 || maxKeys > 1000 {
		maxKeys = 1000 // Use AWS default/max
	}

	input := &s3.ListObjectsV2Input{
		Bucket:    aws.String(bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String(delimiter),
		MaxKeys:   aws.Int32(maxKeys),
	}

	// Only set continuation token if provided
	if nextToken != "" {
		input.ContinuationToken = aws.String(nextToken)
	}

	output, err := s.core.S3Client.ListObjectsV2(ctx, input)
	if err != nil {
		s.core.Logger.Error().
			Err(err).
			Str("bucket", bucket).
			Str("prefix", prefix).
			Msg("Failed to list objects")
		return nil, err
	}

	response := &models.ListObjectsResponse{
		Objects:     make([]models.ObjectInfo, 0, len(output.Contents)+len(output.CommonPrefixes)),
		IsTruncated: *output.IsTruncated,
		PageSize:    int(maxKeys),
	}

	// Process CommonPrefixes (folders)
	for _, prefix := range output.CommonPrefixes {
		response.Objects = append(response.Objects, models.ObjectInfo{
			Key:      aws.ToString(prefix.Prefix),
			IsFolder: true,
			Type:     "folder",
			Size:     0,
		})
	}

	// Process Contents (files)
	for _, obj := range output.Contents {
		key := aws.ToString(obj.Key)

		// Skip the current directory marker if it exists
		if key == prefix {
			continue
		}

		// Basic content type detection based on extension
		contentType := ""
		if !strings.HasSuffix(key, "/") {
			contentType = detectContentType(key)
		}

		response.Objects = append(response.Objects, models.ObjectInfo{
			Key:          key,
			IsFolder:     false,
			Type:         "file",
			Size:         aws.ToInt64(obj.Size),
			ContentType:  contentType,
			LastModified: aws.ToTime(obj.LastModified),
			StorageClass: string(obj.StorageClass),
			ETag:         aws.ToString(obj.ETag),
		})
	}

	response.TotalItems = len(response.Objects)
	return response, nil
}

// GetPresignedURL generates a presigned URL for downloading an object
func (s *S3Service) GetPresignedURL(ctx context.Context, bucket, key string, expiresIn int64) (string, error) {
	s.core.Logger.Debug().
		Str("bucket", bucket).
		Str("key", key).
		Int64("expiresIn", expiresIn).
		Msg("Generating presigned URL")

	if expiresIn <= 0 {
		expiresIn = 15 * 60 // Default to 15 minutes
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	presignClient := s.core.S3Presigner
	resp, err := presignClient.PresignGetObject(ctx, input,
		func(opts *s3.PresignOptions) {
			opts.Expires = time.Duration(expiresIn) * time.Second
		})

	if err != nil {
		s.core.Logger.Error().
			Err(err).
			Str("bucket", bucket).
			Str("key", key).
			Msg("Failed to generate presigned URL")
		return "", err
	}

	return resp.URL, nil
}

// detectContentType detects the content type of a file based on its extension
func detectContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".txt":
		return "text/plain"
	case ".html", ".htm":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".pdf":
		return "application/pdf"
	default:
		return "application/octet-stream"
	}
}
