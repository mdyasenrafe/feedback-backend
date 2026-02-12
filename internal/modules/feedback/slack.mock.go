package feedback

import (
	"context"
	"log"
)

// MockSlackClient is a mock implementation of SlackClient for testing/development.
// It logs messages instead of sending them to Slack.
type MockSlackClient struct{}

// NewMockSlackClient creates a new mock Slack client.
func NewMockSlackClient() *MockSlackClient {
	return &MockSlackClient{}
}

// PublishFeedback logs the feedback message instead of sending to Slack.
// Always returns nil for predictable behavior.
func (m *MockSlackClient) PublishFeedback(ctx context.Context, userEmail string, message string) error {
	log.Printf("[MOCK SLACK] User: %s, Message: %s", userEmail, message)
	return nil
}
