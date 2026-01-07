package migrations

import (
	"database/sql"
	"fmt"
	"log"

	"explorer451/internal/config"

	"github.com/jmoiron/sqlx"
)

func V0_1_0(db *sqlx.DB, config *config.Config, log *log.Logger) error {
	log.Print("Running migration v0.1.0")
	schema := []string{
		`CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY
		)`,

		`CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(255) PRIMARY KEY,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			name VARCHAR(255),
			email VARCHAR(255) UNIQUE,
			username VARCHAR(255) UNIQUE NOT NULL,
			password TEXT, -- Hashed password, can be NULL for OIDC-only users
			status VARCHAR(50) DEFAULT 'active', -- active, inactive, suspended
			role VARCHAR(50) DEFAULT 'user', -- user, admin, etc.
			password_login BOOLEAN DEFAULT true, -- Whether password login is allowed
			loggedin_at TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS tokens (
			id VARCHAR(255) PRIMARY KEY,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			user_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			token VARCHAR(255) UNIQUE NOT NULL,
			name VARCHAR(255) NOT NULL,
			expires_at TIMESTAMP,
			last_used_at TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS sessions (
			id TEXT NOT NULL PRIMARY KEY,
			data jsonb DEFAULT '{}'::jsonb NOT NULL,
			created_at timestamp without time zone DEFAULT now() NOT NULL
		)`,

		`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`,
		`CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)`,
		`CREATE INDEX IF NOT EXISTS idx_tokens_user_id ON tokens(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_tokens_token ON tokens(token)`,
		`CREATE INDEX IF NOT EXISTS idx_sessions ON sessions (id, created_at)`,
	}

	// Start a transaction
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure proper rollback handling
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			// We can only log this error since we can't return it
			fmt.Printf("failed to rollback transaction: %v\n", err)
		}
	}()

	// Execute the schema
	for _, query := range schema {
		if _, err := tx.Exec(query); err != nil {
			return fmt.Errorf("Failed to execute schema: %w", err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}

	return nil
}
