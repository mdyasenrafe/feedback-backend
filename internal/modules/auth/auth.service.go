package auth

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Service struct {
	repo        *Repository
	jwtSecret   string
	mailConfig  MailConfig
	deeplinkURL string
}

func NewService(repo *Repository, jwtSecret, deeplinkURL string, mailConfig MailConfig) *Service {
	return &Service{
		repo:        repo,
		jwtSecret:   jwtSecret,
		mailConfig:  mailConfig,
		deeplinkURL: deeplinkURL,
	}
}

// RequestLoginLink handles the login link request flow.
// It normalizes the email, upserts the user, generates a token, stores it, and sends the email.
func (s *Service) RequestLoginLink(ctx context.Context, email string) error {
	// Normalize email
	normalizedEmail := strings.ToLower(strings.TrimSpace(email))
	if normalizedEmail == "" {
		return fmt.Errorf("email is required")
	}

	// Upsert user
	userID, err := s.repo.UpsertUserByEmail(ctx, normalizedEmail)
	if err != nil {
		return fmt.Errorf("failed to upsert user: %w", err)
	}

	// Generate token
	rawToken, err := GenerateToken()
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}

	// Hash token for storage
	tokenHash := HashToken(rawToken)

	// Create login link with 15 minute expiry
	expiresAt := time.Now().Add(15 * time.Minute)
	if err := s.repo.CreateLoginLink(ctx, userID, tokenHash, expiresAt); err != nil {
		return fmt.Errorf("failed to create login link: %w", err)
	}

	// Send email (do NOT log raw token)
	if err := SendLoginLink(s.mailConfig, normalizedEmail, s.deeplinkURL, rawToken); err != nil {
		fmt.Printf("SendLoginLink error: %v\n", err)
		return fmt.Errorf("email_send_failed: %w", err)
	}

	return nil
}

// VerifyLoginLink verifies the token and returns a JWT and user info.
func (s *Service) VerifyLoginLink(ctx context.Context, rawToken string) (*VerifyLoginLinkResponse, error) {
	if rawToken == "" {
		return nil, fmt.Errorf("token is required")
	}

	// Hash the token
	tokenHash := HashToken(rawToken)

	// Atomically consume the login link
	userID, err := s.repo.ConsumeLoginLink(ctx, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("invalid_or_expired_token")
	}

	// Get user email
	email, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Create JWT
	jwt, err := CreateJWT(userID.String(), email, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to create JWT: %w", err)
	}

	return &VerifyLoginLinkResponse{
		AccessToken: jwt,
		User: User{
			ID:    userID.String(),
			Email: email,
		},
	}, nil
}
