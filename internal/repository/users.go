package repository

import (
	"context"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type usersRepo struct {
	queries *dbsqlc.Queries
}

func (r *usersRepo) Create(ctx context.Context, username, passwordHash, pharmacyID string) (*dbsqlc.User, error) {
	var pharmacyUUID pgtype.UUID
	if err := pharmacyUUID.Scan(pharmacyID); err != nil {
		return nil, err
	}

	params := dbsqlc.CreateUserParams{
		Username:     username,
		PasswordHash: []byte(passwordHash),
		PharmacyID:   pharmacyUUID,
	}
	user, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
