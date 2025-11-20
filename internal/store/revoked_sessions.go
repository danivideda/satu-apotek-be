package store

import (
	"context"
	"errors"
	"time"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type RevokedSessionStore struct {
	queries *dbsqlc.Queries
}

func (s *RevokedSessionStore) Create(ctx context.Context, sessionID string) (*dbsqlc.RevokedSession, error) {
	revokedSession, err := s.queries.CreateRevokedSession(
		ctx,
		dbsqlc.CreateRevokedSessionParams{
			SessionID: sessionID,
			Expires: pgtype.Timestamptz{
				Time:  time.Now(),
				Valid: true,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	return &revokedSession, nil
}

func (s *RevokedSessionStore) GetBySessionID(ctx context.Context, sessionID string) (*dbsqlc.RevokedSession, error) {
	revokedSession, err := s.queries.GetRevokedSessionBySessionID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}
	return &revokedSession, nil
}
