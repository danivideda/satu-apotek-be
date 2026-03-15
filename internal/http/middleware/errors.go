package middleware

import "errors"

var (
	ErrBearerInvalid = errors.New("authorization header is not in 'Bearer <token>' format")
	ErrBearerEmpty   = errors.New("token part of the authorization header is empty")
	ErrAuthMissing   = errors.New("authorization header is missing")
	ErrInvalidRole   = errors.New("invalid role inside JWT")
	ErrInvalidCSRFToken = errors.New("invalid CSRF token")
)
