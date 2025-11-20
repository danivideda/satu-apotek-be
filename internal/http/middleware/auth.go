package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/danivideda/satu-apotek-be/internal/http/jwt"
	"github.com/danivideda/satu-apotek-be/internal/http/response"
)

type contextKey string

const AuthClaimsCtx contextKey = "auth_claims"

func Auth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tokenString, err := extractBearerToken(r)
		if err != nil {
			response.Unauthorized(w, r, err)
			return
		}

		claims, err := jwt.ValidateAccessToken(tokenString)
		if err != nil {
			response.Unauthorized(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, AuthClaimsCtx, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

func AuthClaimsFromContext(ctx context.Context) *jwt.AuthClaims {
	claims, ok := ctx.Value(AuthClaimsCtx).(*jwt.AuthClaims)
	if !ok {
		return nil
	}
	return claims
}

// extractBearerToken takes an Authorization header string and returns the token
// string and an error if the format is incorrect.
func extractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrAuthMissing
	}

	// The standard format is "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", ErrBearerInvalid
	}

	// Trim any potential extra whitespace from the token part
	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", ErrBearerEmpty
	}

	return token, nil
}
