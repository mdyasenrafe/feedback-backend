package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"feedback/internal/shared/httpx"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HandleRequestLoginLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpx.WriteError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		return
	}

	var req RequestLoginLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid_json")
		return
	}

	if err := h.service.RequestLoginLink(r.Context(), req.Email); err != nil {
		// Check if it's an email send failure
		if strings.Contains(err.Error(), "email_send_failed") {
			httpx.WriteError(w, http.StatusInternalServerError, "email_send_failed")
			return
		}
		// Generic error (don't leak details)
		httpx.WriteError(w, http.StatusInternalServerError, "internal_error")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, RequestLoginLinkResponse{OK: true})
}

// HandleVerifyLoginLink handles POST /auth/login-link/verify
func (h *Handler) HandleVerifyLoginLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpx.WriteError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		return
	}

	var req VerifyLoginLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid_json")
		return
	}

	resp, err := h.service.VerifyLoginLink(r.Context(), req.Token)
	if err != nil {
		// Check if it's an invalid/expired token
		if strings.Contains(err.Error(), "invalid_or_expired_token") {
			httpx.WriteError(w, http.StatusUnauthorized, "invalid_or_expired_token")
			return
		}
		// Generic error
		httpx.WriteError(w, http.StatusInternalServerError, "internal_error")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, resp)
}
