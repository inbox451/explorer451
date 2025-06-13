package models

import "time"

// Bucket represents an S3 bucket
type Bucket struct {
	Name         string    `json:"name"`
	CreationDate time.Time `json:"creationDate"`
}

// ObjectInfo represents an S3 object or prefix (folder)
type ObjectInfo struct {
	Key          string    `json:"key"`
	Size         int64     `json:"size"`
	IsFolder     bool      `json:"isFolder"`
	Type         string    `json:"type"`
	ContentType  string    `json:"contentType,omitempty"`
	LastModified time.Time `json:"lastModified"`
	StorageClass string    `json:"storageClass"`
	ETag         string    `json:"etag"`
}

// ListObjectsResponse is the response for listing objects in a bucket
type ListObjectsResponse struct {
	Objects     []ObjectInfo `json:"objects"`
	IsTruncated bool         `json:"isTruncated"`
	ItemsInPage int          `json:"itemsInPage"`
	PageSize    int          `json:"pageSize"`
}

// CreateFolderRequest represents the request body for creating a folder
type CreateFolderRequest struct {
	Key  string `json:"key" validate:"required"`
	Type string `json:"type" validate:"required,eq=folder"`
}

// PresignedPostURLRequest represents the request body for generating a presigned POST URL
type PresignedPostURLRequest struct {
	Key              string `json:"key" validate:"required"`
	ContentType      string `json:"contentType" validate:"required"`
	ExpiresInSeconds int64  `json:"expiresInSeconds,omitempty"`
	MaxSizeBytes     int64  `json:"maxSizeBytes,omitempty"`
}

// PresignedPostURLResponse represents the response for generating a presigned POST URL
type PresignedPostURLResponse struct {
	URL    string            `json:"url"`
	Fields map[string]string `json:"fields"`
}

// BucketDetail represents detailed information about an S3 bucket
type BucketDetail struct {
	Name         string    `json:"name"`
	Region       string    `json:"region"`
	CreationDate time.Time `json:"creationDate"`
}

// ObjectMetadata represents detailed metadata for an S3 object
type ObjectMetadata struct {
	Key                  string            `json:"key"`
	ContentType          string            `json:"contentType"`
	ContentLength        int64             `json:"contentLength"`
	ETag                 string            `json:"etag"`
	LastModified         time.Time         `json:"lastModified"`
	StorageClass         string            `json:"storageClass"`
	UserMetadata         map[string]string `json:"userMetadata,omitempty"`
	ServerSideEncryption string            `json:"serverSideEncryption,omitempty"`
	VersionId            string            `json:"versionId,omitempty"`
}
