package feedback

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Create inserts a new feedback record and returns it.
func (r *Repository) Create(ctx context.Context, userID uuid.UUID, message string) (*Feedback, error) {
	query := `
		INSERT INTO feedback (user_id, message)
		VALUES ($1, $2)
		RETURNING id, user_id, message, created_at
	`

	var f Feedback
	var id, uid uuid.UUID
	err := r.pool.QueryRow(ctx, query, userID, message).Scan(
		&id, &uid, &f.Message, &f.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create feedback: %w", err)
	}

	f.ID = id.String()
	f.UserID = uid.String()
	return &f, nil
}
