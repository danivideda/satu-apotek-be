package repository

import (
	"context"

	"github.com/danivideda/satu-apotek-be/internal/config"
	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ownersRepo struct {
	db         *pgxpool.Pool
	queries    *dbsqlc.Queries
	authConfig config.AuthConfig
}

func (r *ownersRepo) Create(ctx context.Context, username, email, passwordHash string) (ownerID int64, err error) {
	owner, err := r.queries.CreateOwner(ctx, dbsqlc.CreateOwnerParams{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return
	}

	return owner.ID, nil
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
		if isNotFoundError(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &owner, nil
}

func (r *ownersRepo) GetByEmail(ctx context.Context, email string) (*dbsqlc.GetOwnerByEmailRow, error) {
	owner, err := r.queries.GetOwnerByEmail(ctx, email)
	if err != nil {
		if isNotFoundError(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &owner, nil
}
