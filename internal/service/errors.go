package service

import "errors"

var (
	ErrMalformedCSRFToken = errors.New("CSRF token is malformed")
)