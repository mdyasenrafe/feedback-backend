package feedback

import "time"

type CreateFeedbackRequest struct {
	Message string `json:"message"`
}

type Feedback struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
