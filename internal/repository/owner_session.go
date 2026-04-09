package repository

import (
	"context"
	"errors"
	"time"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type ownerSessionsRepo struct {
	queries *dbsqlc.Queries
}

func (r *ownerSessionsRepo) Create(ctx context.Context, ownerID int64, exp time.Time) (*dbsqlc.OwnerSession, error) {
	ownerSession, err := r.queries.CreateOwnerSession(ctx, dbsqlc.CreateOwnerSessionParams{
		OwnerID:   ownerID,
		ExpiresAt: pgtype.Timestamptz{Time: exp, Valid: true},
	},
	)
	if err != nil {
		return nil, err
	}

	return &ownerSession, nil
}

func (r *ownerSessionsRepo) Update(ctx context.Context, sessionID string, exp time.Time) (*dbsqlc.OwnerSession, error) {
	var sessionUUID pgtype.UUID
	if err := sessionUUID.Scan(sessionID); err != nil {
		return nil, err
	}

	ownerSession, err := r.queries.UpdateOwnerSession(ctx, dbsqlc.UpdateOwnerSessionParams{
		SessionID: sessionUUID,
		ExpiresAt: pgtype.Timestamptz{Time: exp, Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}

	return &ownerSession, nil
}

func (r *ownerSessionsRepo) Get(ctx context.Context, sessionID string) (*dbsqlc.OwnerSession, error) {
	var sessionUUID pgtype.UUID
	if err := sessionUUID.Scan(sessionID); err != nil {
		return nil, err
	}

	ownerSession, err := r.queries.GetOwnerSession(ctx, sessionUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}
	return &ownerSession, nil
}

func (r *ownerSessionsRepo) Delete(ctx context.Context, sessionID string) (*dbsqlc.OwnerSession, error) {
	var sessionUUID pgtype.UUID
	if err := sessionUUID.Scan(sessionID); err != nil {
		return nil, err
	}
	deletedOwnerSession, err := r.queries.DeleteOwnerSession(ctx, sessionUUID)
	if err != nil {
		return nil, err
	}
	return &deletedOwnerSession, nil
}

func (r *ownerSessionsRepo) DeleteExpired(ctx context.Context) (*[]dbsqlc.OwnerSession, error) {
	items, err := r.queries.DeleteExpiredOwnerSessions(ctx)
	if err != nil {
		return nil, err
	}
	return &items, nil
}