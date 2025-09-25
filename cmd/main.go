package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"explorer451/internal/api"
	"explorer451/internal/aws"
	"explorer451/internal/config"
	"explorer451/internal/core"
	"explorer451/internal/logger"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Setup logger
	log := logger.New(cfg.Log.Level, cfg.Log.Format)

	// Create context that listens for signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Load AWS configuration
	awsCfg, err := aws.LoadConfig(ctx, &cfg.AWS)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load AWS configuration")
	}

	// Determine if we're using LocalStack
	isLocal := cfg.AWS.EndpointURL != ""

	// Create S3 client
	s3Client := aws.NewS3Client(awsCfg, isLocal)
	s3Presigner := aws.NewS3Presigner(awsCfg, isLocal)

	// Initialize database connection (optional)
	var db *sqlx.DB
	if cfg.Database.URL != "" {
		log.Info().Msg("Initializing database connection")
		db, err = sqlx.Connect("postgres", cfg.Database.URL)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to connect to database, continuing without authentication")
		} else {
			// Configure database connection pool
			db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
			db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
			db.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

			// Test the connection
			if err := db.Ping(); err != nil {
				log.Warn().Err(err).Msg("Database ping failed, continuing without authentication")
				db.Close()
				db = nil
			} else {
				log.Info().Msg("Database connection established")
			}
		}
	}

	// Initialize core service
	core := core.NewCore(cfg, log, s3Client, s3Presigner, db)

	// Setup and start HTTP server
	server, err := api.NewServer(core)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create server")
	}
	go func() {
		if err := server.Start(cfg.Server.Address); err != nil {
			log.Error().Err(err).Msg("Server error")
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()
	log.Info().Msg("Shutdown signal received")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal().Err(err).Msg("Server shutdown failed")
	}

	log.Info().Msg("Server gracefully stopped")
}
