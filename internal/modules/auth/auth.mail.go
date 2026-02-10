package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/resend/resend-go/v2"
)

type MailConfig struct {
	APIKey string
	From   string
}

// SendLoginLink sends an email with the magic link to the user via Resend (HTTPS).
func SendLoginLink(cfg MailConfig, toEmail, deeplinkURL, rawToken string) error {
	if cfg.APIKey == "" {
		return fmt.Errorf("RESEND_API_KEY is required")
	}
	if cfg.From == "" {
		return fmt.Errorf("EMAIL_FROM is required")
	}

	link := fmt.Sprintf("%s?token=%s", deeplinkURL, rawToken)
	body := fmt.Sprintf(
		"Click the link below to log in:\n\n%s\n\nThis link expires in 15 minutes.",
		link,
	)

	client := resend.NewClient(cfg.APIKey)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	params := &resend.SendEmailRequest{
		From:    cfg.From,
		To:      []string{toEmail},
		Subject: "Your login link for FeedbackApp",
		Text:    body,
	}

	_, err := client.Emails.SendWithContext(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to send email via resend: %w", err)
	}

	return nil
}
