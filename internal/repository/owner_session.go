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

func (s *ownerSessionsRepo) Create(ctx context.Context, ownerID int64, exp time.Time) (*dbsqlc.OwnerSession, error) {
	ownerSession, err := s.queries.CreateOwnerSession(ctx, dbsqlc.CreateOwnerSessionParams{
		OwnerID:   ownerID,
		ExpiresAt: pgtype.Timestamptz{Time: exp, Valid: true},
	},
	)
	if err != nil {
		return nil, err
	}

	return &ownerSession, nil
}

func (s *ownerSessionsRepo) Update(ctx context.Context, sessionID string, exp time.Time) (*dbsqlc.OwnerSession, error) {
	var sessionUUID pgtype.UUID
	if err := sessionUUID.Scan(sessionID); err != nil {
		return nil, err
	}

	ownerSession, err := s.queries.UpdateOwnerSession(ctx, dbsqlc.UpdateOwnerSessionParams{
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

func (s *ownerSessionsRepo) Get(ctx context.Context, sessionID string) (*dbsqlc.OwnerSession, error) {
	var sessionUUID pgtype.UUID
	if err := sessionUUID.Scan(sessionID); err != nil {
		return nil, err
	}

	ownerSession, err := s.queries.GetOwnerSession(ctx, sessionUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}
	return &ownerSession, nil
}

func (s *ownerSessionsRepo) Delete(ctx context.Context, sessionID string) (*dbsqlc.OwnerSession, error) {
	var sessionUUID pgtype.UUID
	if err := sessionUUID.Scan(sessionID); err != nil {
		return nil, err
	}
	deletedOwnerSession, err := s.queries.DeleteOwnerSession(ctx, sessionUUID)
	if err != nil {
		return nil, err
	}
	return &deletedOwnerSession, nil
}

func (s *ownerSessionsRepo) DeleteExpired(ctx context.Context) (*[]dbsqlc.OwnerSession, error) {
	items, err := s.queries.DeleteExpiredOwnerSessions(ctx)
	if err != nil {
		return nil, err
	}
	return &items, nil
}