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

// SendLoginLink sends an email with the magic link to the user via Resend.
func SendLoginLink(cfg MailConfig, toEmail, deeplinkURL, rawToken string) error {
	if cfg.APIKey == "" {
		return fmt.Errorf("RESEND_API_KEY is required")
	}
	if cfg.From == "" {
		return fmt.Errorf("EMAIL_FROM is required")
	}

	// deeplinkURL MUST be HTTPS, e.g. https://api.yourdomain.com/auth/deeplink
	link := fmt.Sprintf("%s?token=%s", deeplinkURL, rawToken)

	textBody := fmt.Sprintf(
		"Click the link below to log in:\n\n%s\n\nThis link expires in 15 minutes.",
		link,
	)

	htmlBody := fmt.Sprintf(`
		<div style="font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Arial,sans-serif;line-height:1.4;">
			<p>Click the button below to log in:</p>
			<p>
				<a href="%s" style="display:inline-block;padding:12px 16px;border:1px solid #ccc;border-radius:10px;text-decoration:none;">
					Log in to FeedbackApp
				</a>
			</p>
			<p style="color:#666;">This link expires in 15 minutes.</p>
			<p style="color:#666;">If the button doesn’t work, copy and paste this URL into your browser:</p>
			<p><code>%s</code></p>
		</div>
	`, link, link)

	client := resend.NewClient(cfg.APIKey)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	params := &resend.SendEmailRequest{
		From:    cfg.From,
		To:      []string{toEmail},
		Subject: "Your login link for FeedbackApp",
		Text:    textBody,
		Html:    htmlBody,
	}

	_, err := client.Emails.SendWithContext(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to send email via resend: %w", err)
	}

	return nil
}
