package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
		if strings.Contains(err.Error(), "email_send_failed") {
			httpx.WriteError(w, http.StatusInternalServerError, "email_send_failed")
			return
		}
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
		if strings.Contains(err.Error(), "invalid_or_expired_token") {
			httpx.WriteError(w, http.StatusUnauthorized, "invalid_or_expired_token")
			return
		}
		httpx.WriteError(w, http.StatusInternalServerError, "internal_error")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) HandleDeeplink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpx.WriteError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		return
	}

	rawToken := r.URL.Query().Get("token")
	if rawToken == "" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("missing token"))
		return
	}

	target := "feedbackapp://auth?token=" + url.QueryEscape(rawToken)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// Small HTML that:
	// - tries to open immediately
	// - provides a button fallback
	// Works better than a raw 302 in many in-app browsers.
	_, _ = fmt.Fprintf(w, `<!doctype html>
<html>
<head>
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Open FeedbackApp</title>
</head>
<body style="font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Arial,sans-serif;padding:24px;line-height:1.4;">
  <h2>Opening FeedbackApp…</h2>
  <p>If nothing happens, tap the button below.</p>
  <p>
    <a href="%s"
       style="display:inline-block;padding:12px 16px;border:1px solid #ccc;border-radius:10px;text-decoration:none;">
      Open FeedbackApp
    </a>
  </p>
  <p style="color:#666;margin-top:24px;">This login link expires in 15 minutes.</p>
  <script>
    // Try to open immediately (some clients require a user gesture; button remains as fallback).
    window.location.href = %q;
  </script>
</body>
</html>`, target, target)
}
