package middleware

import (
	"context"
	"curriculum-tracker/utils"
	"net/http"
	"strings"
)

type contextKey string

const UserIDKey contextKey = "userID"

func Auth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Let OPTIONS requests through
			if r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.WriteError(w, http.StatusUnauthorized, "Authorization header required")
				return
			}

			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
				utils.WriteError(w, http.StatusUnauthorized, "Invalid authorization header format")
				return
			}

			claims, err := utils.ValidateToken(bearerToken[1], jwtSecret)
			if err != nil {
				utils.WriteError(w, http.StatusUnauthorized, "Invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(UserIDKey).(int)
	return userID, ok
}
