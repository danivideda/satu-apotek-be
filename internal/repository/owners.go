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
	type txQueryResult struct {
		OwnerID        int64
		OwnerSessionID string
	}
	withTxQuery := func(qtx *dbsqlc.Queries) (any, error) {
		// create new owner
		owner, err := qtx.CreateOwner(ctx, dbsqlc.CreateOwnerParams{
			Username:     username,
			Email:        email,
			PasswordHash: passwordHash,
		})
		if err != nil {
			return nil, err
		}

		// create new session
		ttl, err := time.ParseDuration(ownerSessionTTL)
		if err != nil {
			return nil, err
		}
		ownerSession, err := qtx.CreateOwnerSession(ctx, dbsqlc.CreateOwnerSessionParams{
			OwnerID:   owner.ID,
			ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(ttl), Valid: true},
		})
		if err != nil {
			return nil, err
		}

		result := txQueryResult{
			OwnerID:        owner.ID,
			OwnerSessionID: ownerSession.ID.String(),
		}
		return result, nil
	}

	qtxResult, err := runInTx(ctx, r.db, r.queries, withTxQuery)
	if err != nil {
		return 0, "", err
	}

	owner := qtxResult.(txQueryResult)

	return owner.OwnerID, owner.OwnerSessionID, nil
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
