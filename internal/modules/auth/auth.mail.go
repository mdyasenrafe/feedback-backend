package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// MailConfig is used by SendLoginLink (Mailgun).
type MailConfig struct {
	APIKey  string // MAILGUN_API_KEY
	Domain  string // MAILGUN_DOMAIN, e.g. sandboxxxxx.mailgun.org
	BaseURL string // MAILGUN_BASE_URL, e.g. https://api.mailgun.net (or EU base)
	From    string // EMAIL_FROM, e.g. "FeedbackApp <postmaster@sandboxxxxx.mailgun.org>"
}

// SendLoginLink sends an email with the magic link to the user via Mailgun.
func SendLoginLink(cfg MailConfig, toEmail, deeplinkURL, rawToken string) error {
	if cfg.APIKey == "" || cfg.Domain == "" || cfg.BaseURL == "" || cfg.From == "" {
		return fmt.Errorf("mailgun config missing")
	}
	if toEmail == "" {
		return fmt.Errorf("toEmail is required")
	}
	if deeplinkURL == "" {
		return fmt.Errorf("deeplinkURL is required")
	}
	if rawToken == "" {
		return fmt.Errorf("rawToken is required")
	}

	// deeplinkURL should be an HTTPS (or http in local dev) endpoint that serves /auth/deeplink
	// e.g. https://your-api.onrender.com/auth/deeplink
	link := fmt.Sprintf("%s?token=%s", strings.TrimRight(deeplinkURL, "/"), url.QueryEscape(rawToken))

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

	form := url.Values{}
	form.Set("from", cfg.From)
	form.Set("to", toEmail)
	form.Set("subject", "Your login link for FeedbackApp")
	form.Set("text", textBody)
	form.Set("html", htmlBody)

	endpoint := strings.TrimRight(cfg.BaseURL, "/") + "/v3/" + cfg.Domain + "/messages"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("mailgun request build failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Basic auth: username "api", password API key
	authHeader := base64.StdEncoding.EncodeToString([]byte("api:" + cfg.APIKey))
	req.Header.Set("Authorization", "Basic "+authHeader)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("mailgun send failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("mailgun send failed: status=%d body=%s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}
