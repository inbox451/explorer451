package storage

import (
	"context"
	"fmt"

	"explorer451/internal/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Repository interface {
	// User operations
	ListUsers(ctx context.Context, limit, offset int) ([]*models.User, int, error)
	GetUser(ctx context.Context, id string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, userId string) error

	// Token operations
	ListTokensByUser(ctx context.Context, userID string, limit, offset int) ([]*models.Token, int, error)
	GetTokenByUser(ctx context.Context, userID string, tokenID string) (*models.Token, error)
	CreateToken(ctx context.Context, token *models.Token) error
	DeleteToken(ctx context.Context, tokenID string) error
	GetTokenByValue(ctx context.Context, tokenValue string) (*models.Token, error)
	UpdateTokenLastUsed(ctx context.Context, tokenID string) error
	PruneExpiredTokens(ctx context.Context) (int64, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

// User operations
func (r *repository) ListUsers(ctx context.Context, limit, offset int) ([]*models.User, int, error) {
	var users []*models.User
	var total int

	// Get total count
	err := r.db.GetContext(ctx, &total, "SELECT COUNT(*) FROM users")
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Get users with pagination
	query := `SELECT id, created_at, updated_at, name, email, username, password, status, role, password_login, loggedin_at 
			  FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err = r.db.SelectContext(ctx, &users, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, total, nil
}

func (r *repository) GetUser(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	query := `SELECT id, created_at, updated_at, name, email, username, password, status, role, password_login, loggedin_at 
			  FROM users WHERE id = $1`
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *repository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	query := `SELECT id, created_at, updated_at, name, email, username, password, status, role, password_login, loggedin_at 
			  FROM users WHERE username = $1`
	err := r.db.GetContext(ctx, &user, query, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return &user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, created_at, updated_at, name, email, username, password, status, role, password_login, loggedin_at 
			  FROM users WHERE email = $1`
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

func (r *repository) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (id, name, email, username, password, status, role, password_login, loggedin_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.Username, user.Password, user.Status, user.Role, user.PasswordLogin, user.LoggedinAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *repository) UpdateUser(ctx context.Context, user *models.User) error {
	query := `UPDATE users SET name = $2, email = $3, username = $4, password = $5, status = $6, role = $7, password_login = $8, loggedin_at = $9, updated_at = NOW() 
			  WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.Username, user.Password, user.Status, user.Role, user.PasswordLogin, user.LoggedinAt)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *repository) DeleteUser(ctx context.Context, userId string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, userId)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// Token operations
func (r *repository) ListTokensByUser(ctx context.Context, userID string, limit, offset int) ([]*models.Token, int, error) {
	var tokens []*models.Token
	var total int

	// Get total count
	err := r.db.GetContext(ctx, &total, "SELECT COUNT(*) FROM tokens WHERE user_id = $1", userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tokens: %w", err)
	}

	// Get tokens with pagination
	query := `SELECT id, created_at, updated_at, user_id, token, name, expires_at, last_used_at 
			  FROM tokens WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	err = r.db.SelectContext(ctx, &tokens, query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list tokens: %w", err)
	}

	return tokens, total, nil
}

func (r *repository) GetTokenByUser(ctx context.Context, userID string, tokenID string) (*models.Token, error) {
	var token models.Token
	query := `SELECT id, created_at, updated_at, user_id, token, name, expires_at, last_used_at 
			  FROM tokens WHERE user_id = $1 AND id = $2`
	err := r.db.GetContext(ctx, &token, query, userID, tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	return &token, nil
}

func (r *repository) CreateToken(ctx context.Context, token *models.Token) error {
	query := `INSERT INTO tokens (id, user_id, token, name, expires_at, last_used_at) 
			  VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, token.ID, token.UserID, token.Token, token.Name, token.ExpiresAt, token.LastUsedAt)
	if err != nil {
		return fmt.Errorf("failed to create token: %w", err)
	}
	return nil
}

func (r *repository) DeleteToken(ctx context.Context, tokenID string) error {
	query := `DELETE FROM tokens WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, tokenID)
	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}
	return nil
}

func (r *repository) GetTokenByValue(ctx context.Context, tokenValue string) (*models.Token, error) {
	var token models.Token
	query := `SELECT id, created_at, updated_at, user_id, token, name, expires_at, last_used_at 
			  FROM tokens WHERE token = $1`
	err := r.db.GetContext(ctx, &token, query, tokenValue)
	if err != nil {
		return nil, fmt.Errorf("failed to get token by value: %w", err)
	}
	return &token, nil
}

func (r *repository) UpdateTokenLastUsed(ctx context.Context, tokenID string) error {
	query := `UPDATE tokens SET last_used_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, tokenID)
	if err != nil {
		return fmt.Errorf("failed to update token last used: %w", err)
	}
	return nil
}

func (r *repository) PruneExpiredTokens(ctx context.Context) (int64, error) {
	query := `DELETE FROM tokens WHERE expires_at IS NOT NULL AND expires_at < NOW()`
	result, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to prune expired tokens: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}
	return rowsAffected, nil
}
