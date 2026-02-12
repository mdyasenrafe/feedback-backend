package feedback

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
)

type Service struct {
	repo        *Repository
	slackClient SlackClient
}

func NewService(repo *Repository, slackClient SlackClient) *Service {
	return &Service{
		repo:        repo,
		slackClient: slackClient,
	}
}

func (s *Service) CreateFeedback(ctx context.Context, userID uuid.UUID, userEmail, message string) (*Feedback, error) {
	normalizedMessage := strings.TrimSpace(message)

	if normalizedMessage == "" {
		return nil, fmt.Errorf("message_required")
	}

	// Persist feedback (DB is source of truth)
	feedback, err := s.repo.Create(ctx, userID, normalizedMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to create feedback: %w", err)
	}

	// Publish to Slack (best effort - don't fail request if this fails)
	if err := s.slackClient.PublishFeedback(ctx, userEmail, normalizedMessage); err != nil {
		log.Printf("Slack publish failed for feedback %s: %v", feedback.ID, err)
		// Continue - feedback was stored successfully
	}

	return feedback, nil
}
