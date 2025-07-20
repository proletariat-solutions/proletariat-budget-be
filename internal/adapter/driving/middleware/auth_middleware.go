package middleware

import (
	"context"
	"net/http"
	"strings"

	"ghorkov32/proletariat-budget-be/internal/core/usecase"
)

type AuthMiddleware struct {
	authUseCase *usecase.AuthUseCase
}

const (
	USER_CONTEXT_KEY = "user"
)

func NewAuthMiddleware(authUseCase *usecase.AuthUseCase) *AuthMiddleware {
	return &AuthMiddleware{
		authUseCase: authUseCase,
	}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(
					w,
					"Unauthorized",
					http.StatusUnauthorized,
				)

				return
			}

			// Check if it's a Bearer token
			parts := strings.Split(
				authHeader,
				" ",
			)
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(
					w,
					"Invalid authorization format",
					http.StatusUnauthorized,
				)

				return
			}

			token := parts[1]

			// Validate token
			user, err := m.authUseCase.ValidateToken(
				r.Context(),
				token,
			)
			if err != nil {
				http.Error(
					w,
					"Unauthorized",
					http.StatusUnauthorized,
				)

				return
			}

			// Add user to context
			ctx := context.WithValue(
				r.Context(),
				USER_CONTEXT_KEY,
				user,
			)

			// Call the next handler with the authenticated context
			next.ServeHTTP(
				w,
				r.WithContext(ctx),
			)
		},
	)
}
