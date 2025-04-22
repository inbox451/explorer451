package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// LoadConfig loads AWS configuration using the default credential chain
func LoadConfig(ctx context.Context, region string) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithRetryMaxAttempts(3),
	)
}

// NewS3Client creates a new S3 client
func NewS3Client(cfg aws.Config) *s3.Client {
	return s3.NewFromConfig(cfg)
}

// NewS3Presigner creates a new S3 presigner client
func NewS3Presigner(cfg aws.Config) *s3.PresignClient {
	return s3.NewPresignClient(s3.NewFromConfig(cfg))
}
