package middlewares

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const (
	userKey  contextKey = "user_id"
	adminKey contextKey = "is_admin"
)

type AuthnMiddleware struct {
}

func NewAuthnMiddleware() *AuthnMiddleware {
	return &AuthnMiddleware{}
}

func (authn *AuthnMiddleware) RequireSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		role := r.Header.Get("X-User-Role")

		if userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if _, err := uuid.Parse(userID); err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, userID)
		ctx = context.WithValue(ctx, adminKey, role == "admin")

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (authn *AuthnMiddleware) RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAdmin, ok := r.Context().Value(adminKey).(bool)
		if !ok || !isAdmin {
			http.Error(w, "Not allowed", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetUserIDFromContext(ctx context.Context) *uuid.UUID {
	if usr, ok := ctx.Value(userKey).(uuid.UUID); ok {
		return &usr
	}
	return nil
}
