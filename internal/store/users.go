package store

import (
	"context"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserStore struct {
	queries *dbsqlc.Queries
}

func (s *UserStore) Create(ctx context.Context, username, passwordHash, pharmacyID string) (*dbsqlc.User, error) {
	var pharmacyUUID pgtype.UUID
	if err := pharmacyUUID.Scan(pharmacyID); err != nil {
		return nil, err
	}

	params := dbsqlc.CreateUserParams{
		Username:     username,
		PasswordHash: []byte(passwordHash),
		PharmacyID:   pharmacyUUID,
	}
	user, err := s.queries.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
