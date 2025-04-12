package core

import (
	"explorer451/internal/config"
	"explorer451/internal/logger"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Core holds the application's core components and services
type Core struct {
	Config      *config.Config
	Logger      *logger.Logger
	S3Client    *s3.Client
	S3Presigner *s3.PresignClient
	S3Service   *S3Service
}

// NewCore creates a new Core instance with all dependencies
func NewCore(
	cfg *config.Config,
	logger *logger.Logger,
	s3Client *s3.Client,
	s3Presigner *s3.PresignClient,
) *Core {
	core := &Core{
		Config:      cfg,
		Logger:      logger,
		S3Client:    s3Client,
		S3Presigner: s3Presigner,
	}

	// Initialize services
	core.S3Service = NewS3Service(core)

	return core
}
