package repository

import (
	"context"
	"time"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ownersRepo struct {
	db      *pgxpool.Pool
	queries *dbsqlc.Queries
}

func (r *ownersRepo) Create(ctx context.Context, username, email, passwordHash string) (ownerID int64, ownerSessionID string, err error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, "", err
	}
	defer tx.Rollback(ctx)
	qtx := r.queries.WithTx(tx)
	// create new owner
	owner, err := qtx.CreateOwner(ctx, dbsqlc.CreateOwnerParams{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return 0, "", err
	}

	// create new session
	ttl, err := time.ParseDuration(ownerSessionTTL)
	if err != nil {
		return 0, "", err
	}
	ownerSession, err := qtx.CreateOwnerSession(ctx, dbsqlc.CreateOwnerSessionParams{
		OwnerID:   owner.ID,
		ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(ttl), Valid: true},
	})
	if err != nil {
		return 0, "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, "", err
	}

	return owner.ID, ownerSession.ID.String(), nil
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
