package aws

import (
	"context"

	appconfig "explorer451/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// LoadConfig loads AWS configuration using the default credential chain
// with support for custom endpoints (LocalStack)
func LoadConfig(ctx context.Context, cfg *appconfig.AWSConfig) (aws.Config, error) {
	var opts []func(*config.LoadOptions) error

	opts = append(opts,
		config.WithRegion(cfg.Region),
		config.WithRetryMaxAttempts(3),
	)

	// If custom endpoint is provided, configure for LocalStack
	if cfg.EndpointURL != "" {
		// LocalStack uses static credentials
		opts = append(opts,
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "test")),
			config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL:               cfg.EndpointURL,
						HostnameImmutable: true,
					}, nil
				})),
		)
	}

	return config.LoadDefaultConfig(ctx, opts...)
}

// NewS3Client creates a new S3 client
func NewS3Client(cfg aws.Config, usePathStyle bool) *s3.Client {
	opts := []func(*s3.Options){
		func(o *s3.Options) {
			if usePathStyle {
				o.UsePathStyle = true
			}
		},
	}
	return s3.NewFromConfig(cfg, opts...)
}

// NewS3Presigner creates a new S3 presigner client
func NewS3Presigner(cfg aws.Config, usePathStyle bool) *s3.PresignClient {
	return s3.NewPresignClient(NewS3Client(cfg, usePathStyle))
}
