package middleware

import (
	"avito-kitchen-service/internal/service/auth"
	"avito-kitchen-service/internal/transport/http/handler"
	"avito-kitchen-service/internal/transport/http/middleware/constants"
	"context"
	"net/http"
	"strings"
)

func AuthMiddleware(jwtService *auth.JWTService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				handler.WriteError(w, "INVALID_REQUEST", "invalid request", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := jwtService.ParseToken(tokenStr)
			if err != nil {
				handler.WriteError(w, "INVALID_REQUEST", "invalid request", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), constants.UserIDKey, claims.UserID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
