package core

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"explorer451/internal/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
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

// GetBucketDetails retrieves detailed information about a bucket including its region
func (s *S3Service) GetBucketDetails(ctx context.Context, bucketName string) (*models.BucketDetail, error) {
	s.core.Logger.Debug().Str("bucket", bucketName).Msg("Getting bucket details")

	// Get bucket location/region
	locationResp, err := s.core.S3Client.GetBucketLocation(ctx, &s3.GetBucketLocationInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		s.core.Logger.Error().Err(err).Str("bucket", bucketName).Msg("Failed to get bucket location")
		return nil, err
	}

	// The location constraint can be empty for us-east-1
	region := string(locationResp.LocationConstraint)
	if region == "" {
		region = "us-east-1"
	}

	// Get bucket creation date from ListBuckets
	bucketsResp, err := s.core.S3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		s.core.Logger.Error().Err(err).Msg("Failed to list buckets")
		return nil, err
	}

	var creationDate time.Time
	for _, bucket := range bucketsResp.Buckets {
		if aws.ToString(bucket.Name) == bucketName {
			creationDate = aws.ToTime(bucket.CreationDate)
			break
		}
	}

	return &models.BucketDetail{
		Name:         bucketName,
		Region:       region,
		CreationDate: creationDate,
	}, nil
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

	response.ItemsInPage = len(response.Objects)
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

// GetObjectMetadata retrieves detailed metadata for an S3 object
func (s *S3Service) GetObjectMetadata(ctx context.Context, bucket, key string) (*models.ObjectMetadata, error) {
	s.core.Logger.Debug().
		Str("bucket", bucket).
		Str("key", key).
		Msg("Getting object metadata")

	output, err := s.core.S3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		s.core.Logger.Error().
			Err(err).
			Str("bucket", bucket).
			Str("key", key).
			Msg("Failed to get object metadata")
		return nil, err
	}

	metadata := &models.ObjectMetadata{
		Key:           key,
		ContentType:   aws.ToString(output.ContentType),
		ContentLength: aws.ToInt64(output.ContentLength),
		ETag:          aws.ToString(output.ETag),
		LastModified:  aws.ToTime(output.LastModified),
		StorageClass:  string(output.StorageClass),
		UserMetadata:  output.Metadata,
		VersionId:     aws.ToString(output.VersionId),
	}

	// Add server-side encryption info if present
	if output.ServerSideEncryption != "" {
		metadata.ServerSideEncryption = string(output.ServerSideEncryption)
	}

	return metadata, nil
}

// GeneratePresignedPostURL generates a presigned POST URL for uploading objects
func (s *S3Service) GeneratePresignedPostURL(ctx context.Context, bucket, key, contentType string, expiresIn time.Duration, maxSize int64) (*models.PresignedPostURLResponse, error) {
	s.core.Logger.Debug().
		Str("bucket", bucket).
		Str("key", key).
		Str("contentType", contentType).
		Dur("expiresIn", expiresIn).
		Int64("maxSize", maxSize).
		Msg("Generating presigned POST URL")

	if expiresIn <= 0 {
		expiresIn = 15 * time.Minute // Default to 15 minutes
	}

	// Set default max size if not specified (10MB)
	if maxSize <= 0 {
		maxSize = 10 * 1024 * 1024 // 10MB
	}

	// Create presigned POST policy
	resp, err := s.core.S3Presigner.PresignPostObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, func(opts *s3.PresignPostOptions) {
		opts.Expires = expiresIn
		opts.Conditions = append(opts.Conditions,
			// Restrict content type
			[]interface{}{"eq", "$Content-Type", contentType},
			// Restrict content length
			[]interface{}{"content-length-range", 0, maxSize},
		)
	})
	if err != nil {
		s.core.Logger.Error().
			Err(err).
			Str("bucket", bucket).
			Str("key", key).
			Msg("Failed to generate presigned POST URL")
		return nil, err
	}

	return &models.PresignedPostURLResponse{
		URL:    resp.URL,
		Fields: resp.Values,
	}, nil
}

// DeleteObject deletes a single object from S3
func (s *S3Service) DeleteObject(ctx context.Context, bucket, key string) error {
	s.core.Logger.Debug().
		Str("bucket", bucket).
		Str("key", key).
		Msg("Deleting object")

	_, err := s.core.S3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		s.core.Logger.Error().
			Err(err).
			Str("bucket", bucket).
			Str("key", key).
			Msg("Failed to delete object")
		return err
	}

	s.core.Logger.Info().
		Str("bucket", bucket).
		Str("key", key).
		Msg("Successfully deleted object")

	return nil
}

// DeleteObjectsByPrefix deletes all objects with the given prefix (folder deletion)
func (s *S3Service) DeleteObjectsByPrefix(ctx context.Context, bucket, prefix string) error {
	s.core.Logger.Debug().
		Str("bucket", bucket).
		Str("prefix", prefix).
		Msg("Deleting objects by prefix")

	// First, list all objects with the prefix
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	}

	var objectsToDelete []s3Types.ObjectIdentifier
	paginator := s3.NewListObjectsV2Paginator(s.core.S3Client, input)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			s.core.Logger.Error().
				Err(err).
				Str("bucket", bucket).
				Str("prefix", prefix).
				Msg("Failed to list objects for deletion")
			return err
		}

		for _, obj := range page.Contents {
			objectsToDelete = append(objectsToDelete, s3Types.ObjectIdentifier{
				Key: obj.Key,
			})
		}
	}

	if len(objectsToDelete) == 0 {
		s.core.Logger.Info().
			Str("bucket", bucket).
			Str("prefix", prefix).
			Msg("No objects found with prefix, nothing to delete")
		return nil
	}

	// Delete objects in batches of 1000 (AWS limit)
	batchSize := 1000
	for i := 0; i < len(objectsToDelete); i += batchSize {
		end := i + batchSize
		if end > len(objectsToDelete) {
			end = len(objectsToDelete)
		}

		batch := objectsToDelete[i:end]
		_, err := s.core.S3Client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
			Bucket: aws.String(bucket),
			Delete: &s3Types.Delete{
				Objects: batch,
				Quiet:   aws.Bool(false), // Get detailed response
			},
		})
		if err != nil {
			s.core.Logger.Error().
				Err(err).
				Str("bucket", bucket).
				Str("prefix", prefix).
				Int("batchStart", i).
				Int("batchEnd", end).
				Msg("Failed to delete batch of objects")
			return err
		}

		s.core.Logger.Info().
			Str("bucket", bucket).
			Str("prefix", prefix).
			Int("count", len(batch)).
			Msg("Successfully deleted batch of objects")
	}

	s.core.Logger.Info().
		Str("bucket", bucket).
		Str("prefix", prefix).
		Int("totalDeleted", len(objectsToDelete)).
		Msg("Successfully deleted all objects with prefix")

	return nil
}

// CreateFolder creates a "folder" in S3 by creating a zero-byte object with a trailing slash
func (s *S3Service) CreateFolder(ctx context.Context, bucket, key string) error {
	s.core.Logger.Debug().
		Str("bucket", bucket).
		Str("key", key).
		Msg("Creating folder")

	// Ensure key ends with '/'
	if !strings.HasSuffix(key, "/") {
		key = key + "/"
	}

	_, err := s.core.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   strings.NewReader(""), // Empty body for folder marker
	})
	if err != nil {
		s.core.Logger.Error().
			Err(err).
			Str("bucket", bucket).
			Str("key", key).
			Msg("Failed to create folder")
		return err
	}

	s.core.Logger.Info().
		Str("bucket", bucket).
		Str("key", key).
		Msg("Successfully created folder")

	return nil
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
