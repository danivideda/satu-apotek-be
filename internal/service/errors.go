package service

import "errors"

var (
	ErrMalformedCSRFToken = errors.New("CSRF token is malformed")
	ErrInvalidCSRFToken = errors.New("invalid CSRF token")
)