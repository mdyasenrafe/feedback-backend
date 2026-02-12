package feedback

import "context"

type SlackClient interface {
	PublishFeedback(ctx context.Context, userEmail string, message string) error
}
