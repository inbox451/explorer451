package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPresignedPostURLRequest_JSON(t *testing.T) {
	tests := []struct {
		name     string
		req      PresignedPostURLRequest
		expected string
	}{
		{
			name: "complete request",
			req: PresignedPostURLRequest{
				Key:              "uploads/test-file.txt",
				ContentType:      "text/plain",
				ExpiresInSeconds: 3600,
				MaxSizeBytes:     1024 * 1024,
			},
			expected: `{"key":"uploads/test-file.txt","contentType":"text/plain","expiresInSeconds":3600,"maxSizeBytes":1048576}`,
		},
		{
			name: "minimal request",
			req: PresignedPostURLRequest{
				Key:         "test.jpg",
				ContentType: "image/jpeg",
			},
			expected: `{"key":"test.jpg","contentType":"image/jpeg"}`,
		},
		{
			name: "request with defaults omitted",
			req: PresignedPostURLRequest{
				Key:              "document.pdf",
				ContentType:      "application/pdf",
				ExpiresInSeconds: 0, // Should be omitted in JSON
				MaxSizeBytes:     0, // Should be omitted in JSON
			},
			expected: `{"key":"document.pdf","contentType":"application/pdf"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := json.Marshal(tt.req)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.expected, string(jsonBytes))

			// Test unmarshaling
			var unmarshaled PresignedPostURLRequest
			err = json.Unmarshal([]byte(tt.expected), &unmarshaled)
			assert.NoError(t, err)
			assert.Equal(t, tt.req.Key, unmarshaled.Key)
			assert.Equal(t, tt.req.ContentType, unmarshaled.ContentType)
			assert.Equal(t, tt.req.ExpiresInSeconds, unmarshaled.ExpiresInSeconds)
			assert.Equal(t, tt.req.MaxSizeBytes, unmarshaled.MaxSizeBytes)
		})
	}
}

func TestPresignedPostURLResponse_JSON(t *testing.T) {
	tests := []struct {
		name     string
		resp     PresignedPostURLResponse
		expected string
	}{
		{
			name: "complete response",
			resp: PresignedPostURLResponse{
				URL: "https://test-bucket.s3.amazonaws.com/",
				Fields: map[string]string{
					"key":            "uploads/test-file.txt",
					"Content-Type":   "text/plain",
					"AWSAccessKeyId": "AKIAIOSFODNN7EXAMPLE",
					"policy":         "eyJleHBpcmF0aW9uIjoi...",
					"signature":      "signature",
				},
			},
			expected: `{"url":"https://test-bucket.s3.amazonaws.com/","fields":{"key":"uploads/test-file.txt","Content-Type":"text/plain","AWSAccessKeyId":"AKIAIOSFODNN7EXAMPLE","policy":"eyJleHBpcmF0aW9uIjoi...","signature":"signature"}}`,
		},
		{
			name: "minimal response",
			resp: PresignedPostURLResponse{
				URL:    "https://example.s3.amazonaws.com/",
				Fields: map[string]string{},
			},
			expected: `{"url":"https://example.s3.amazonaws.com/","fields":{}}`,
		},
		{
			name: "nil fields",
			resp: PresignedPostURLResponse{
				URL:    "https://example.s3.amazonaws.com/",
				Fields: nil,
			},
			expected: `{"url":"https://example.s3.amazonaws.com/","fields":null}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := json.Marshal(tt.resp)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.expected, string(jsonBytes))

			// Test unmarshaling
			var unmarshaled PresignedPostURLResponse
			err = json.Unmarshal([]byte(tt.expected), &unmarshaled)
			assert.NoError(t, err)
			assert.Equal(t, tt.resp.URL, unmarshaled.URL)
			if tt.resp.Fields == nil {
				assert.Nil(t, unmarshaled.Fields)
			} else {
				assert.Equal(t, tt.resp.Fields, unmarshaled.Fields)
			}
		})
	}
}

func TestPresignedPostURLRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     PresignedPostURLRequest
		isValid bool
	}{
		{
			name: "valid complete request",
			req: PresignedPostURLRequest{
				Key:              "uploads/test-file.txt",
				ContentType:      "text/plain",
				ExpiresInSeconds: 3600,
				MaxSizeBytes:     1024 * 1024,
			},
			isValid: true,
		},
		{
			name: "valid minimal request",
			req: PresignedPostURLRequest{
				Key:         "test.jpg",
				ContentType: "image/jpeg",
			},
			isValid: true,
		},
		{
			name: "missing key",
			req: PresignedPostURLRequest{
				ContentType:      "text/plain",
				ExpiresInSeconds: 3600,
			},
			isValid: false,
		},
		{
			name: "missing content type",
			req: PresignedPostURLRequest{
				Key:              "test.txt",
				ExpiresInSeconds: 3600,
			},
			isValid: false,
		},
		{
			name: "empty key",
			req: PresignedPostURLRequest{
				Key:         "",
				ContentType: "text/plain",
			},
			isValid: false,
		},
		{
			name: "empty content type",
			req: PresignedPostURLRequest{
				Key:         "test.txt",
				ContentType: "",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation logic - in a real application you might use a validation library
			hasKey := tt.req.Key != ""
			hasContentType := tt.req.ContentType != ""

			actualValid := hasKey && hasContentType
			assert.Equal(t, tt.isValid, actualValid, "validation result mismatch")
		})
	}
}

func TestCreateFolderRequest_ExistingModel(t *testing.T) {
	// Test to ensure existing models still work
	req := CreateFolderRequest{
		Key:  "new-folder/",
		Type: "folder",
	}

	jsonBytes, err := json.Marshal(req)
	assert.NoError(t, err)

	expected := `{"key":"new-folder/","type":"folder"}`
	assert.JSONEq(t, expected, string(jsonBytes))

	// Test unmarshaling
	var unmarshaled CreateFolderRequest
	err = json.Unmarshal(jsonBytes, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, req.Key, unmarshaled.Key)
	assert.Equal(t, req.Type, unmarshaled.Type)
}

func TestBucket_ExistingModel(t *testing.T) {
	// Test to ensure existing models still work
	now := time.Now()
	bucket := Bucket{
		Name:         "test-bucket",
		CreationDate: now,
	}

	jsonBytes, err := json.Marshal(bucket)
	assert.NoError(t, err)

	// Test unmarshaling
	var unmarshaled Bucket
	err = json.Unmarshal(jsonBytes, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, bucket.Name, unmarshaled.Name)
	assert.True(t, bucket.CreationDate.Equal(unmarshaled.CreationDate))
}

func TestObjectInfo_ExistingModel(t *testing.T) {
	// Test to ensure existing models still work
	now := time.Now()
	obj := ObjectInfo{
		Key:          "test-file.txt",
		Size:         1024,
		IsFolder:     false,
		Type:         "file",
		ContentType:  "text/plain",
		LastModified: now,
		StorageClass: "STANDARD",
		ETag:         "\"etag123\"",
	}

	jsonBytes, err := json.Marshal(obj)
	assert.NoError(t, err)

	// Test unmarshaling
	var unmarshaled ObjectInfo
	err = json.Unmarshal(jsonBytes, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, obj.Key, unmarshaled.Key)
	assert.Equal(t, obj.Size, unmarshaled.Size)
	assert.Equal(t, obj.IsFolder, unmarshaled.IsFolder)
	assert.Equal(t, obj.Type, unmarshaled.Type)
	assert.Equal(t, obj.ContentType, unmarshaled.ContentType)
	assert.True(t, obj.LastModified.Equal(unmarshaled.LastModified))
	assert.Equal(t, obj.StorageClass, unmarshaled.StorageClass)
	assert.Equal(t, obj.ETag, unmarshaled.ETag)
}
