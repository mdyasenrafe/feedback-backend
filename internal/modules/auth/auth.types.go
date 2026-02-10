package auth

import (
	"time"

	"github.com/google/uuid"
)

// Request/Response types

type RequestLoginLinkRequest struct {
	Email string `json:"email"`
}

type RequestLoginLinkResponse struct {
	OK bool `json:"ok"`
}

type VerifyLoginLinkRequest struct {
	Token string `json:"token"`
}

type VerifyLoginLinkResponse struct {
	AccessToken string `json:"accessToken"`
	User        User   `json:"user"`
}

// Domain types

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type LoginLink struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}
