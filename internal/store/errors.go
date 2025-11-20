package store

import "errors"

var (
	ErrNotFound = errors.New("item not found in database")
)
