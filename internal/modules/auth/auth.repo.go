package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// UpsertUserByEmail creates a user if they don't exist, or returns the existing user ID.
func (r *Repository) UpsertUserByEmail(ctx context.Context, email string) (uuid.UUID, error) {
	var userID uuid.UUID
	query := `
		INSERT INTO users (email)
		VALUES ($1)
		ON CONFLICT (email) DO UPDATE SET email = EXCLUDED.email
		RETURNING id
	`
	err := r.pool.QueryRow(ctx, query, email).Scan(&userID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to upsert user: %w", err)
	}
	return userID, nil
}

// CreateLoginLink inserts a new login link with the given token hash and expiry.
func (r *Repository) CreateLoginLink(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	query := `
		INSERT INTO login_links (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`
	_, err := r.pool.Exec(ctx, query, userID, tokenHash, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to create login link: %w", err)
	}
	return nil
}

// ConsumeLoginLink atomically marks a login link as used and returns the user ID.
// Returns an error if the token is invalid, expired, or already used.
func (r *Repository) ConsumeLoginLink(ctx context.Context, tokenHash string) (uuid.UUID, error) {
	var userID uuid.UUID
	query := `
		UPDATE login_links
		SET used_at = now()
		WHERE token_hash = $1
		  AND used_at IS NULL
		  AND expires_at > now()
		RETURNING user_id
	`
	err := r.pool.QueryRow(ctx, query, tokenHash).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, fmt.Errorf("invalid or expired token")
		}
		return uuid.Nil, fmt.Errorf("failed to consume login link: %w", err)
	}
	return userID, nil
}

// GetUserByID retrieves a user's email by their ID.
func (r *Repository) GetUserByID(ctx context.Context, userID uuid.UUID) (string, error) {
	var email string
	query := `SELECT email FROM users WHERE id = $1`
	err := r.pool.QueryRow(ctx, query, userID).Scan(&email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("user not found")
		}
		return "", fmt.Errorf("failed to get user: %w", err)
	}
	return email, nil
}
