# Explorer451

This is a simple file explorer that allows you to navigate through the files and directories of S3.


## Local Development Setup

### Requirements

- Docker and Docker Compose installed on your machine.
- Install `curl` for making HTTP requests.
- Install `awslocal` via `brew install awscli-local` for local AWS S3 interactions.
- Install `awscli` for AWS CLI commands.
- (Optional) Install `jq` for parsing JSON responses in shell scripts.
- (Optional) Install `make` for easier command execution.


### Prepopulate AWS CLI for LocalStack

Set up AWS CLI with LocalStack default region

```bash
export AWS_DEFAULT_REGION=us-east-1
```

Create a sample bucket and upload a sample file (README.md) to it:

```bash
awslocal s3api create-bucket --bucket sample-bucket

awslocal s3api put-object \
  --bucket sample-bucket \
  --key README.md \
  --body README.md
```

## Explorer451 API Documentation

## Base URL

```
http://localhost:8080/api
```

## Authentication

The API supports optional authentication. If authentication is enabled, you'll need to include authentication headers or use the auth endpoints first.

## API

### Login

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "username",
    "password": "password"
  }'
```

**Login Response Example:**
```json
{
  "user": {
    "id": "123...",
    "username": "john",
    "name": "John Doe",
    "email": "user@example.com",
    "role": "user"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires": "2023-12-31T23:59:59Z"
}
```

### Authentication Methods

The API supports 3 authentication methods:

#### 1. Bearer Token Authentication

Extract the token from the login response and use it in the Authorization header:

```bash
# Get the token from login response
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "your-username",
    "password": "your-password"
  }' | jq -r '.token')

# Use the token in subsequent requests
curl -X GET http://localhost:8080/api/auth/profile \
  -H "Authorization: Bearer $TOKEN"
```

#### 2. Session Cookie Authentication

The login endpoint also sets a session cookie that can be used for authentication:

```bash
# Login and save cookies to file
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "your-username",
    "password": "your-password"
  }' \
  -c cookies.txt

# Use the saved cookies for subsequent requests
curl -X GET http://localhost:8080/api/auth/profile \
  -b cookies.txt
```


#### 3. OIDC Authentication (if enabled)

```bash
# Initiate OIDC login
curl -X GET http://localhost:8080/api/auth/oidc/login

# Handle OIDC callback (usually done by browser)
curl -X GET "http://localhost:8080/api/auth/oidc/callback?code=authorization-code&state=state-value"
```

### Get Profile (requires authentication)
```bash
# Using Bearer Token
curl -X GET http://localhost:8080/api/auth/profile \
  -H "Authorization: Bearer $TOKEN"

# Using Session Cookie
curl -X GET http://localhost:8080/api/auth/profile \
  -b cookies.txt
```

### Logout
```bash
# Using Bearer Token
curl -X POST http://localhost:8080/api/auth/logout \
  -H "Authorization: Bearer your-token"

# Using Session Cookie
curl -X POST http://localhost:8080/api/auth/logout \
  -b cookies.txt
```

### Refresh Token
```bash
# Refresh your bearer token (requires authentication)
curl -X POST http://localhost:8080/api/auth/refresh \
  -H "Authorization: Bearer your-current-token"
```

**Refresh Token Response:**
```json
{
  "user": {
    "id": "user-123",
    "username": "your-username",
    "name": "John Doe",
    "email": "user@example.com",
    "role": "user"
  },
  "token": "new-bearer-token-here",
  "expires": "2023-12-31T23:59:59Z"
}
```

### Revoke Token
```bash
# Revoke your current bearer token (requires authentication)
curl -X POST http://localhost:8080/api/auth/revoke \
  -H "Authorization: Bearer $TOKEN"
```

**Revoke Token Response:**
```json
{
  "message": "Token revoked successfully"
}
```

## Bucket Operations

### List All Buckets
```bash
# Without authentication (if auth is disabled)
curl -X GET http://localhost:8080/api/buckets

# With Bearer Token
curl -X GET http://localhost:8080/api/buckets \
  -H "Authorization: Bearer $TOKEN"

# With Session Cookie
curl -X GET http://localhost:8080/api/buckets \
  -b cookies.txt
```

### Get Bucket Details
```bash
# With Bearer Token
curl -X GET http://localhost:8080/api/buckets/sample-bucket \
  -H "Authorization: Bearer $TOKEN"

# With Session Cookie
curl -X GET http://localhost:8080/api/buckets/sample-bucket \
  -b cookies.txt
```

### Create Bucket
```bash
# With Bearer Token
curl -X PUT http://localhost:8080/api/buckets/new-api-bucket \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "region": "us-west-2"
  }'

# With Session Cookie
curl -X PUT http://localhost:8080/api/buckets/new-api-bucket \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "region": "us-west-2"
  }'
```

### Delete Bucket
```bash
# With Bearer Token
curl -X DELETE http://localhost:8080/api/buckets/new-api-bucket \
  -H "Authorization: Bearer $TOKEN"

# With Session Cookie
curl -X DELETE http://localhost:8080/api/buckets/new-api-bucket \
  -b cookies.txt
```

## Object Operations

### List Objects in Bucket
```bash
# List all objects (with Bearer Token)
curl -X GET http://localhost:8080/api/buckets/sample-bucket/objects \
  -H "Authorization: Bearer $TOKEN"

# List objects with prefix (folder)
curl -X GET "http://localhost:8080/api/buckets/sample-bucket/objects?prefix=documents/" \
  -H "Authorization: Bearer $TOKEN"

# List objects with pagination
curl -X GET "http://localhost:8080/api/buckets/sample-bucket/objects?maxKeys=50&nextToken=next-token" \
  -H "Authorization: Bearer $TOKEN"

# List objects with delimiter
curl -X GET "http://localhost:8080/api/buckets/sample-bucket/objects?delimiter=/" \
  -H "Authorization: Bearer $TOKEN"

# Using Session Cookie
curl -X GET http://localhost:8080/api/buckets/sample-bucket/objects \
  -b cookies.txt
```

### Download Object
```bash
# Download object directly (with Bearer Token)
curl -X GET http://localhost:8080/api/buckets/sample-bucket/objects/README.md \
  -H "Authorization: Bearer $TOKEN" \
  -o README.md.txt

# Download with specific headers
curl -X GET http://localhost:8080/api/buckets/sample-bucket/objects/README.md \
  -H "Authorization: Bearer $TOKEN" \
  -H "Accept: application/octet-stream" \
  -o README.md.txt

# Using Session Cookie
curl -X GET http://localhost:8080/api/buckets/sample-bucket/objects/path/to/file.txt \
  -b cookies.txt \
  -o README.md.txt
```

### Upload Object
```bash
# Upload a file (with Bearer Token)
curl -X PUT http://localhost:8080/api/buckets/sample-bucket/objects/path/to/docker-compose.yml \
  -H "Content-Type: text/plain" \
  -H "Authorization: Bearer $TOKEN" \
  --data-binary @docker-compose.yml

# Upload with custom metadata
curl -X PUT http://localhost:8080/api/buckets/my-bucket/objects/path/to/uploaded-file.txt \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -H "x-amz-meta-custom-field: custom-value" \
  --data-binary @local-file.json

# Using Session Cookie
curl -X PUT http://localhost:8080/api/buckets/my-bucket/objects/path/to/uploaded-file.txt \
  -H "Content-Type: text/plain" \
  -b cookies.txt \
  --data-binary @local-file.txt
```

### Get Object Metadata
```bash
# With Bearer Token
curl -X HEAD http://localhost:8080/api/buckets/my-bucket/objects/path/to/file.txt \
  -H "Authorization: Bearer $TOKEN" \
  -v

# Using Session Cookie
curl -X HEAD http://localhost:8080/api/buckets/my-bucket/objects/path/to/file.txt \
  -b cookies.txt \
  -v
```

### Delete Object
```bash
# Delete single object (with Bearer Token)
curl -X DELETE http://localhost:8080/api/buckets/my-bucket/objects/path/to/file.txt \
  -H "Authorization: Bearer $TOKEN"

# Delete folder recursively
curl -X DELETE "http://localhost:8080/api/buckets/my-bucket/objects/path/to/folder/?recursive=true" \
  -H "Authorization: Bearer $TOKEN"

# Using Session Cookie
curl -X DELETE http://localhost:8080/api/buckets/my-bucket/objects/path/to/file.txt \
  -b cookies.txt
```

### Copy Object
```bash
# Copy object within same bucket (with Bearer Token)
curl -X POST http://localhost:8080/api/buckets/my-bucket/objects/path/to/copied-file.txt \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "sourceKey": "path/to/original-file.txt"
  }'

# Copy object from different bucket
curl -X POST http://localhost:8080/api/buckets/my-bucket/objects/path/to/copied-file.txt \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "sourceKey": "path/to/original-file.txt",
    "sourceBucket": "source-bucket",
    "contentType": "text/plain",
    "storageClass": "STANDARD"
  }'

# Using Session Cookie
curl -X POST http://localhost:8080/api/buckets/my-bucket/objects/path/to/copied-file.txt \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "sourceKey": "path/to/original-file.txt"
  }'
```

## Presigned URL Operations

### Generate Download URL
```bash
# Generate download URL with default expiration (15 minutes) - with Bearer Token
curl -X POST http://localhost:8080/api/buckets/my-bucket/presigned/download/path/to/file.txt \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{}'

# Generate download URL with custom expiration (1 hour)
curl -X POST http://localhost:8080/api/buckets/my-bucket/presigned/download/path/to/file.txt \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "expiresIn": 3600
  }'

# Using Session Cookie
curl -X POST http://localhost:8080/api/buckets/my-bucket/presigned/download/path/to/file.txt \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{}'
```

### Generate Upload URL
```bash
# Generate upload URL (with Bearer Token)
curl -X POST http://localhost:8080/api/buckets/my-bucket/presigned/upload/path/to/upload-file.txt \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "contentType": "text/plain",
    "expiresIn": 3600,
    "maxSizeBytes": 10485760
  }'

# Using Session Cookie
curl -X POST http://localhost:8080/api/buckets/my-bucket/presigned/upload/path/to/upload-file.txt \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "contentType": "text/plain",
    "expiresIn": 3600,
    "maxSizeBytes": 10485760
  }'
```

## Multipart Upload Operations

### Initiate Multipart Upload
```bash
# With Bearer Token
curl -X POST "http://localhost:8080/api/buckets/my-bucket/uploads?key=path/to/large-file.zip" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "contentType": "application/zip",
    "storageClass": "STANDARD",
    "serverSideEncryption": "AES256"
  }'

# Using Session Cookie
curl -X POST "http://localhost:8080/api/buckets/my-bucket/uploads?key=path/to/large-file.zip" \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "contentType": "application/zip",
    "storageClass": "STANDARD",
    "serverSideEncryption": "AES256"
  }'
```

### Upload Part
```bash
# Upload part 1 (with Bearer Token)
curl -X PUT "http://localhost:8080/api/buckets/my-bucket/uploads/upload-id-123/parts/1?key=path/to/large-file.zip" \
  -H "Content-Type: application/octet-stream" \
  -H "Authorization: Bearer $TOKEN" \
  --data-binary @part1.bin

# Upload part 2
curl -X PUT "http://localhost:8080/api/buckets/my-bucket/uploads/upload-id-123/parts/2?key=path/to/large-file.zip" \
  -H "Content-Type: application/octet-stream" \
  -H "Authorization: Bearer $TOKEN" \
  --data-binary @part2.bin

# Using Session Cookie
curl -X PUT "http://localhost:8080/api/buckets/my-bucket/uploads/upload-id-123/parts/1?key=path/to/large-file.zip" \
  -H "Content-Type: application/octet-stream" \
  -b cookies.txt \
  --data-binary @part1.bin
```

### Complete Multipart Upload
```bash
# With Bearer Token
curl -X POST "http://localhost:8080/api/buckets/my-bucket/uploads/upload-id-123/complete?key=path/to/large-file.zip" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "parts": [
      {
        "etag": "etag-from-part-1",
        "partNumber": 1
      },
      {
        "etag": "etag-from-part-2",
        "partNumber": 2
      }
    ]
  }'

# Using Session Cookie
curl -X POST "http://localhost:8080/api/buckets/my-bucket/uploads/upload-id-123/complete?key=path/to/large-file.zip" \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "parts": [
      {
        "etag": "etag-from-part-1",
        "partNumber": 1
      },
      {
        "etag": "etag-from-part-2",
        "partNumber": 2
      }
    ]
  }'
```

### Abort Multipart Upload
```bash
# With Bearer Token
curl -X DELETE "http://localhost:8080/api/buckets/my-bucket/uploads/upload-id-123?key=path/to/large-file.zip" \
  -H "Authorization: Bearer $TOKEN"

# Using Session Cookie
curl -X DELETE "http://localhost:8080/api/buckets/my-bucket/uploads/upload-id-123?key=path/to/large-file.zip" \
  -b cookies.txt
```

## Batch Operations

### Batch Delete Objects
```bash
# With Bearer Token
curl -X POST http://localhost:8080/api/buckets/my-bucket/objects:batchDelete \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "objects": [
      {
        "key": "path/to/file1.txt"
      },
      {
        "key": "path/to/file2.txt",
        "versionId": "version-id-if-versioned"
      },
      {
        "key": "path/to/file3.txt"
      }
    ]
  }'

# Using Session Cookie
curl -X POST http://localhost:8080/api/buckets/my-bucket/objects:batchDelete \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "objects": [
      {
        "key": "path/to/file1.txt"
      },
      {
        "key": "path/to/file2.txt",
        "versionId": "version-id-if-versioned"
      },
      {
        "key": "path/to/file3.txt"
      }
    ]
  }'
```
