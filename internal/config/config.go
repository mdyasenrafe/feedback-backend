package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port             int
	DatabaseURL      string
	JWTSecret        string
	AppDeeplinkURL   string
	SMTPHost         string
	SMTPPort         int
	SMTPUser         string
	SMTPPass         string
	SMTPFrom         string
}

// Load reads and validates all required environment variables.
// It fails fast with a clear error message if any required variable is missing.
func Load() (*Config, error) {
	cfg := &Config{}

	// PORT (default 8080)
	portStr := os.Getenv("PORT")
	if portStr == "" {
		cfg.Port = 8080
	} else {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("invalid PORT: %w", err)
		}
		cfg.Port = port
	}

	// DATABASE_URL (required)
	cfg.DatabaseURL = os.Getenv("DATABASE_URL")
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	// JWT_SECRET (required)
	cfg.JWTSecret = os.Getenv("JWT_SECRET")
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	// APP_DEEPLINK_URL (required)
	cfg.AppDeeplinkURL = os.Getenv("APP_DEEPLINK_URL")
	if cfg.AppDeeplinkURL == "" {
		return nil, fmt.Errorf("APP_DEEPLINK_URL is required")
	}

	// SMTP_HOST (required)
	cfg.SMTPHost = os.Getenv("SMTP_HOST")
	if cfg.SMTPHost == "" {
		return nil, fmt.Errorf("SMTP_HOST is required")
	}

	// SMTP_PORT (required)
	smtpPortStr := os.Getenv("SMTP_PORT")
	if smtpPortStr == "" {
		return nil, fmt.Errorf("SMTP_PORT is required")
	}
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return nil, fmt.Errorf("invalid SMTP_PORT: %w", err)
	}
	cfg.SMTPPort = smtpPort

	// SMTP_USER (required)
	cfg.SMTPUser = os.Getenv("SMTP_USER")
	if cfg.SMTPUser == "" {
		return nil, fmt.Errorf("SMTP_USER is required")
	}

	// SMTP_PASS (required)
	cfg.SMTPPass = os.Getenv("SMTP_PASS")
	if cfg.SMTPPass == "" {
		return nil, fmt.Errorf("SMTP_PASS is required")
	}

	// SMTP_FROM (required)
	cfg.SMTPFrom = os.Getenv("SMTP_FROM")
	if cfg.SMTPFrom == "" {
		return nil, fmt.Errorf("SMTP_FROM is required")
	}

	return cfg, nil
}
