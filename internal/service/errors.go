package service

import "errors"

var (
	ErrMalformedCSRFToken = errors.New("CSRF token is malformed")
	ErrInvalidCSRFToken   = errors.New("invalid CSRF token")
)

var (
	ErrInvalidCredentials       = errors.New("invalid credentials")
	ErrInvalidOwnerSessionCache = errors.New("type assertion failed, ownerID is not int64")
)
