package feedback

import (
	"context"
	"database/sql"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Create(ctx context.Context, userID int64, message string) (*Feedback, error) {
	const q = `
		INSERT INTO feedback (user_id, message)
		VALUES ($1, $2)
		RETURNING id, user_id, message, created_at;
	`

	var f Feedback
	err := r.db.QueryRowContext(ctx, q, userID, message).Scan(
		&f.ID, &f.UserID, &f.Message, &f.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &f, nil
}
