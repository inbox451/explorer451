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
	"github.com/spf13/pflag"
)

var (
	installFlag    = pflag.Bool("install", false, "Initialize the database schema")
	upgradeFlag    = pflag.Bool("upgrade", false, "Upgrade the database schema")
	yesFlag        = pflag.Bool("yes", false, "Skip confirmation prompts")
	idempotentFlag = pflag.Bool("idempotent", false, "Allow idempotent installs")
)

func main() {
	// Parse command line flags
	pflag.Parse()

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

	// Initialize database connection (required for install/upgrade commands)
	var db *sqlx.DB
	if cfg.Database.URL != "" {
		log.Info().Msg("Initializing database connection")
		db, err = sqlx.Connect("postgres", cfg.Database.URL)
		if err != nil {
			if *installFlag || *upgradeFlag {
				log.Fatal().Err(err).Msg("Failed to connect to database for install/upgrade")
			}
			log.Warn().Err(err).Msg("Failed to connect to database, continuing without authentication")
		} else {
			// Configure database connection pool
			db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
			db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
			db.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

			// Test the connection
			if err := db.Ping(); err != nil {
				if *installFlag || *upgradeFlag {
					log.Fatal().Err(err).Msg("Database ping failed for install/upgrade")
				}
				log.Warn().Err(err).Msg("Database ping failed, continuing without authentication")
				db.Close()
				db = nil
			} else {
				log.Info().Msg("Database connection established")
			}
		}
	} else if *installFlag || *upgradeFlag {
		log.Fatal().Msg("Database URL is required for install/upgrade commands")
	}

	// Handle install command
	if *installFlag {
		if db == nil {
			log.Fatal().Msg("Database connection is required for install")
		}
		install(db, cfg, !*yesFlag, *idempotentFlag)
		os.Exit(0)
	}

	// Handle upgrade command
	if *upgradeFlag {
		if db == nil {
			log.Fatal().Msg("Database connection is required for upgrade")
		}
		upgrade(db, cfg, !*yesFlag)
		os.Exit(0)
	}

	// Check if the DB schema is installed (only if database is available)
	if db != nil {
		checkInstall(db)
		checkUpgrade(db)
	}

	// Initialize core service
	core := core.NewCore(cfg, log, s3Client, s3Presigner, db)

	// Setup and start HTTP server
	server, err := api.NewServer(core)
	log.Info().Msg("Starting server... at " + cfg.Server.Address)
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
