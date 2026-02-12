package middleware

import (
	"context"
	"net/http"
	"strings"

	"feedback/internal/shared/httpx"

	"github.com/golang-jwt/jwt/v5"
)

// Context keys for storing user information
type contextKey string

const (
	userIDKey    contextKey = "userID"
	userEmailKey contextKey = "userEmail"
)

// JWTClaims matches the structure from auth.jwt.go
type JWTClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// RequireAuth is a middleware that validates JWT tokens and extracts user identity.
// It expects the Authorization header in the format: "Bearer <token>"
func RequireAuth(jwtSecret string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Extract Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				httpx.WriteError(w, http.StatusUnauthorized, "missing_authorization")
				return
			}

			// Check Bearer prefix
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				httpx.WriteError(w, http.StatusUnauthorized, "invalid_authorization_format")
				return
			}

			tokenString := parts[1]

			// Parse and validate JWT
			token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
				// Verify signing method is HS256
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				httpx.WriteError(w, http.StatusUnauthorized, "invalid_token")
				return
			}

			// Extract claims
			claims, ok := token.Claims.(*JWTClaims)
			if !ok {
				httpx.WriteError(w, http.StatusUnauthorized, "invalid_claims")
				return
			}

			// Store user info in context
			ctx := context.WithValue(r.Context(), userIDKey, claims.Subject)
			ctx = context.WithValue(ctx, userEmailKey, claims.Email)

			// Call next handler with updated context
			next(w, r.WithContext(ctx))
		}
	}
}

// GetAuthUser extracts the authenticated user's ID and email from the request context.
// Returns (userID, email, true) if found, or ("", "", false) if not found.
func GetAuthUser(r *http.Request) (userID string, email string, ok bool) {
	userID, ok1 := r.Context().Value(userIDKey).(string)
	email, ok2 := r.Context().Value(userEmailKey).(string)
	return userID, email, ok1 && ok2
}
