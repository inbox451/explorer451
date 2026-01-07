package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

const (
	// EnvPrefix is the prefix for environment variables
	EnvPrefix = "EXPLORER451_"
)

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	URL             string        `koanf:"url"`
	MaxOpenConns    int           `koanf:"max_open_conns"`
	MaxIdleConns    int           `koanf:"max_idle_conns"`
	ConnMaxLifetime time.Duration `koanf:"conn_max_lifetime"`
}

// OIDCConfig holds OpenID Connect / OAuth2 settings
type OIDCConfig struct {
	Enabled      bool   `koanf:"enabled"`
	ProviderURL  string `koanf:"provider_url"`
	RedirectURL  string `koanf:"redirect_url"`
	ClientID     string `koanf:"client_id"`
	ClientSecret string `koanf:"client_secret"`
}

// Config holds all application configuration
type Config struct {
	Server   ServerConfig   `koanf:"server"`
	AWS      AWSConfig      `koanf:"aws"`
	Log      LogConfig      `koanf:"log"`
	Database DatabaseConfig `koanf:"database"`
	OIDC     OIDCConfig     `koanf:"oidc"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Address       string `koanf:"address"`
	SecureCookies bool   `koanf:"secure_cookies"` // Set to true when using HTTPS
}

// AWSConfig holds AWS specific configuration
type AWSConfig struct {
	Region      string `koanf:"region"`
	EndpointURL string `koanf:"endpoint_url"`
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
		// Convert EXPLORER451_SERVER_ADDRESS to server.address
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

	if cfg.Database.MaxOpenConns == 0 {
		cfg.Database.MaxOpenConns = 25
	}

	if cfg.Database.MaxIdleConns == 0 {
		cfg.Database.MaxIdleConns = 5
	}

	if cfg.Database.ConnMaxLifetime == 0 {
		cfg.Database.ConnMaxLifetime = 5 * time.Minute
	}
}
