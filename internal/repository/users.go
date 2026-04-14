package repository

import (
	"context"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
)

type usersRepo struct {
	queries *dbsqlc.Queries
}

func (r *usersRepo) Create(ctx context.Context, username, passwordHash string, pharmacyID int64) (*dbsqlc.User, error) {
	params := dbsqlc.CreateUserParams{
		Username:     username,
		PasswordHash: passwordHash,
		PharmacyID:   pharmacyID,
	}
	user, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		if isDuplicateError(err) {
			return nil, ErrDuplicateValue
		}
		return nil, err
	}

	return &user, nil
}
