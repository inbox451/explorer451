#!/bin/bash
echo "--- Initializing S3 Buckets ---"

# Set dummy AWS credentials for LocalStack
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_DEFAULT_REGION=us-east-1

# Use AWS CLI with LocalStack endpoint
aws --endpoint-url=http://localhost:4566 s3 mb s3://my-local-bucket
aws --endpoint-url=http://localhost:4566 s3 mb s3://another-local-bucket

echo "--- S3 Initialization Complete ---"