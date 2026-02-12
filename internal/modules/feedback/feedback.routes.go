package feedback

import (
	"net/http"

	"feedback/internal/middleware"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RegisterRoutes registers all feedback routes on the provided mux.
func RegisterRoutes(mux *http.ServeMux, pool *pgxpool.Pool, jwtSecret string) {
	repo := NewRepository(pool)
	slackClient := NewMockSlackClient()
	service := NewService(repo, slackClient)
	handler := NewHandler(service)

	// POST /feedback - requires authentication
	mux.HandleFunc("/feedback", middleware.RequireAuth(jwtSecret)(handler.HandleCreateFeedback))
}
