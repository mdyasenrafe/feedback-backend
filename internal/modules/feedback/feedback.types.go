package feedback

import "time"

// Request/Response types

type CreateFeedbackRequest struct {
	Message string `json:"message"`
}

type CreateFeedbackResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

// Domain types

type Feedback struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
