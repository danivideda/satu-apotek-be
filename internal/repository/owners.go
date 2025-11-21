package repository

import (
	"context"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type OwnersRepo struct {
	queries *dbsqlc.Queries
}

func (r *OwnersRepo) Create(ctx context.Context, username, email, passwordHash string) (*dbsqlc.CreateOwnerRow, error) {
	params := dbsqlc.CreateOwnerParams{
		Username:     username,
		Email:        email,
		PasswordHash: []byte(passwordHash),
	}

	owner, err := r.queries.CreateOwner(ctx, params)
	if err != nil {
		return nil, err
	}
	return &owner, nil
}

func (r *OwnersRepo) GetByID(ctx context.Context, id string) (*dbsqlc.Owner, error) {
	var ownerUUID pgtype.UUID
	if err := ownerUUID.Scan(id); err != nil {
		return nil, err
	}
	owner, err := r.queries.GetOwnerByID(ctx, ownerUUID)
	if err != nil {
		return nil, err
	}
	return &owner, nil
}

func (r *OwnersRepo) GetByUsername(ctx context.Context, username string) (*dbsqlc.GetOwnerByUsernameRow, error) {
	owner, err := r.queries.GetOwnerByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return &owner, nil
}
