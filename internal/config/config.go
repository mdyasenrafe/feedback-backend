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

	ResendAPIKey string
	EmailFrom    string
}

func Load() (Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	cfg := Config{
		Port:           port,
		DatabaseURL:    mustEnv("DATABASE_URL"),
		JWTSecret:      mustEnv("JWT_SECRET"),
		AppDeeplinkURL: mustEnv("APP_DEEPLINK_URL"),

		ResendAPIKey: mustEnv("RESEND_API_KEY"),
		EmailFrom:    mustEnv("EMAIL_FROM"),
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
