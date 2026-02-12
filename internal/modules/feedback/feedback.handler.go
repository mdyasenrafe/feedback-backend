package feedback

import (
	"encoding/json"
	"net/http"
	"strings"

	"feedback/internal/middleware"
	"feedback/internal/shared/httpx"

	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// HandleCreateFeedback handles POST /feedback
func (h *Handler) HandleCreateFeedback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpx.WriteError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		return
	}

	// Extract authenticated user from context (set by middleware)
	userIDStr, userEmail, ok := middleware.GetAuthUser(r)
	if !ok {
		// Should never happen if middleware is working correctly
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse user ID as UUID
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		httpx.WriteError(w, http.StatusUnauthorized, "invalid_user_id")
		return
	}

	// Decode request body
	var req CreateFeedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid_json")
		return
	}

	// Create feedback
	feedback, err := h.service.CreateFeedback(r.Context(), userID, userEmail, req.Message)
	if err != nil {
		// Check for validation errors
		if strings.Contains(err.Error(), "message_required") {
			httpx.WriteError(w, http.StatusBadRequest, "message_required")
			return
		}
		if strings.Contains(err.Error(), "message_too_long") {
			httpx.WriteError(w, http.StatusBadRequest, "message_too_long")
			return
		}
		// Generic error
		httpx.WriteError(w, http.StatusInternalServerError, "internal_error")
		return
	}

	// Return created feedback
	httpx.WriteJSON(w, http.StatusCreated, feedback)
}
