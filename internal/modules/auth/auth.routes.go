package auth

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RegisterRoutes registers all auth routes on the provided mux.
func RegisterRoutes(mux *http.ServeMux, pool *pgxpool.Pool, jwtSecret, deeplinkURL string, mailConfig MailConfig) {
	repo := NewRepository(pool)
	service := NewService(repo, jwtSecret, deeplinkURL, mailConfig)
	handler := NewHandler(service)

	mux.HandleFunc("/auth/login-link", handler.HandleRequestLoginLink)
	mux.HandleFunc("/auth/login-link/verify", handler.HandleVerifyLoginLink)
	mux.HandleFunc("/auth/deeplink", handler.HandleDeeplink)
}
