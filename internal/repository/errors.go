package repository

import "errors"

var (
	ErrNotFound = errors.New("item not found in database")
	ErrSqidsConfigNotSet = errors.New("sqids config not set")
)
