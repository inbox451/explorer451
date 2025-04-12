package config

import (
	"fmt"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

const (
	// EnvPrefix is the prefix for environment variables
	EnvPrefix = "S3EXPLORER_"
)

// Config holds all application configuration
type Config struct {
	Server ServerConfig `koanf:"server"`
	AWS    AWSConfig    `koanf:"aws"`
	Log    LogConfig    `koanf:"log"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Address string `koanf:"address"`
}

// AWSConfig holds AWS specific configuration
type AWSConfig struct {
	Region string `koanf:"region"`
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level  string `koanf:"level"`
	Format string `koanf:"format"`
}

// Load loads configuration from config file and environment variables
func Load() (*Config, error) {
	k := koanf.New(".")

	// Load default configuration
	if err := k.Load(file.Provider("config.yml"), yaml.Parser()); err != nil {
		// Config file is optional, only log error if it exists but can't be loaded
		if !strings.Contains(err.Error(), "no such file or directory") {
			return nil, fmt.Errorf("error loading config file: %w", err)
		}
	}

	// Load environment variables
	callback := func(s string) string {
		// Convert S3EXPLORER_SERVER_ADDRESS to server.address
		path := strings.Replace(strings.ToLower(strings.TrimPrefix(s, EnvPrefix)), "_", ".", -1)
		return path
	}

	if err := k.Load(env.Provider(EnvPrefix, ".", callback), nil); err != nil {
		return nil, fmt.Errorf("error loading environment variables: %w", err)
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	applyDefaults(&cfg)
	return &cfg, nil
}

// applyDefaults sets sensible defaults for empty config values
func applyDefaults(cfg *Config) {
	if cfg.Server.Address == "" {
		cfg.Server.Address = ":8080"
	}

	if cfg.AWS.Region == "" {
		cfg.AWS.Region = "us-east-1"
	}

	if cfg.Log.Level == "" {
		cfg.Log.Level = "info"
	}

	if cfg.Log.Format == "" {
		cfg.Log.Format = "json"
	}
}
