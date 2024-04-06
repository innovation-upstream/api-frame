package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/innovation-upstream/api-frame/provider"
)

// Private key for context. This is to prevent collisions between different context uses
var uidCtxKey = &contextKey{"fb_auth_uid"}

type contextKey struct {
	name string
}

type errorResponse struct {
	Message string `json:"message"`
}

// Authenticate returns a middleware function that parses and validates a firebase session JWT.
// The authenticated user's UID is then injected into the request context.
// Also supports custom firebase tokens
func Authenticate() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("session")
			if err != nil || c == nil {
				msg, err := json.Marshal(errorResponse{Message: "unauthorized"})
				if err != nil {
					panic(err)
				}

				http.Error(w, string(msg), http.StatusUnauthorized)
				return
			}

			firebaseClient := provider.NewFirebaseClient(r.Context())
			claims, err := firebaseClient.Auth().VerifySessionCookie(r.Context(), c.Value)
			if err != nil {
				msg, err := json.Marshal(errorResponse{Message: "unauthorized"})
				if err != nil {
					panic(err)
				}

				http.Error(w, string(msg), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), uidCtxKey, claims.UID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// GetFirebaseUIDFromContext retrieves the Firebase UID from a request context and panics if it is unset.
func GetFirebaseUIDFromContext(ctx context.Context) string {
	raw := ctx.Value(uidCtxKey)
	if raw == nil {
		panic("Failed to get Firebase UID from context.")
	}

	return raw.(string)
}
