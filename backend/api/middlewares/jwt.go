package middlewares

import (
	"context"
	"net/http"

	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/services/jwt"
)

// ContextKey is the key for the context value
type ContextKey string

const (
	// KeyUserID is the id for the context key
	KeyUserID = ContextKey("user-id")
)

// EnrichUserID adds the user id to the context
func EnrichUserID(jwtService jwt.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			//
			// Allow unauthenticated users in
			if header == "" {
				next.ServeHTTP(w, r)
				return
			}

			//validate jwt token
			tokenString := header
			userID, err := jwtService.GetUserIDFromToken(tokenString)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}

			// put it in context
			ctx := context.WithValue(r.Context(), KeyUserID, userID)

			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
