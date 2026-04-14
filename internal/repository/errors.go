package repository

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrNotFound          = errors.New("item not found in database")
	ErrSqidsConfigNotSet = errors.New("sqids config not set")
	ErrDuplicateValue    = errors.New("duplicate value")
)

func isDuplicateError(err error) bool {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return true
			}
		}
		return false
}