package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port           string
	DatabaseURL    string
	JWTSecret      string
	AppDeeplinkURL string

	MailgunAPIKey  string
	MailgunDomain  string
	MailgunBaseURL string
	EmailFrom      string
}

func Load() (Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Default to US API base if not set
	mailgunBaseURL := os.Getenv("MAILGUN_BASE_URL")
	if mailgunBaseURL == "" {
		mailgunBaseURL = "https://api.mailgun.net"
	}

	cfg := Config{
		Port:           port,
		DatabaseURL:    mustEnv("DATABASE_URL"),
		JWTSecret:      mustEnv("JWT_SECRET"),
		AppDeeplinkURL: mustEnv("APP_DEEPLINK_URL"),

		MailgunAPIKey:  mustEnv("MAILGUN_API_KEY"),
		MailgunDomain:  mustEnv("MAILGUN_DOMAIN"), // e.g. sandboxXXXX.mailgun.org
		MailgunBaseURL: mailgunBaseURL,
		EmailFrom:      mustEnv("EMAIL_FROM"),
	}

	return cfg, nil
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("%s is required", key))
	}
	return v
}
