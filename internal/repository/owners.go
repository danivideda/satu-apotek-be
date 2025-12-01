package repository

import (
	"context"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
)

type ownersRepo struct {
	queries *dbsqlc.Queries
}

func (r *ownersRepo) Create(ctx context.Context, username, email, passwordHash string) (*dbsqlc.CreateOwnerRow, error) {
	params := dbsqlc.CreateOwnerParams{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
	}

	owner, err := r.queries.CreateOwner(ctx, params)
	if err != nil {
		return nil, err
	}
	return &owner, nil
}

func (r *ownersRepo) GetByID(ctx context.Context, id int64) (*dbsqlc.Owner, error) {
	owner, err := r.queries.GetOwnerByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &owner, nil
}

func (r *ownersRepo) GetByUsername(ctx context.Context, username string) (*dbsqlc.GetOwnerByUsernameRow, error) {
	owner, err := r.queries.GetOwnerByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return &owner, nil
}
