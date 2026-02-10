package auth

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

type MailConfig struct {
	Host string
	Port int
	User string
	Pass string
	From string
}

// SendLoginLink sends an email with the magic link to the user.
func SendLoginLink(cfg MailConfig, toEmail, deeplinkURL, rawToken string) error {
	link := fmt.Sprintf("%s?token=%s", deeplinkURL, rawToken)

	m := gomail.NewMessage()
	m.SetHeader("From", cfg.From)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Your login link for FeedbackApp")
	m.SetBody("text/plain", fmt.Sprintf("Click the link below to log in:\n\n%s\n\nThis link expires in 15 minutes.", link))

	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.User, cfg.Pass)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
