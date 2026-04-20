package repository

import (
	"context"
	"errors"
	"time"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type userSessionsRepo struct {
	queries *dbsqlc.Queries
}

func (r *userSessionsRepo) Create(ctx context.Context, userID int64, exp time.Time) (*dbsqlc.UserSession, error) {
	userSession, err := r.queries.CreateUserSession(ctx, dbsqlc.CreateUserSessionParams{
		UserID:    userID,
		ExpiresAt: pgtype.Timestamptz{Time: exp, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &userSession, nil
}

func (r *userSessionsRepo) Update(ctx context.Context, sessionID string, exp time.Time) (*dbsqlc.UserSession, error) {
	var sessionUUID pgtype.UUID
	if err := sessionUUID.Scan(sessionID); err != nil {
		return nil, err
	}

	userSession, err := r.queries.UpdateUserSession(ctx, dbsqlc.UpdateUserSessionParams{
		SessionID: sessionUUID,
		ExpiresAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &userSession, nil
}

func (r *userSessionsRepo) Get(ctx context.Context, sessionID string) (*dbsqlc.UserSession, error) {
	var sessionUUID pgtype.UUID
	if err := sessionUUID.Scan(sessionID); err != nil {
		return nil, err
	}

	userSession, err := r.queries.GetUserSession(ctx, sessionUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}
	return &userSession, nil
}

func (r *userSessionsRepo) Delete(ctx context.Context, sessionID string) (*dbsqlc.UserSession, error) {
	var sessionUUID pgtype.UUID
	if err := sessionUUID.Scan(sessionID); err != nil {
		return nil, err
	}
	deletedUserSession, err := r.queries.DeleteUserSession(ctx, sessionUUID)
	if err != nil {
		return nil, err
	}
	return &deletedUserSession, nil
}

func (r *userSessionsRepo) DeleteExpired(ctx context.Context) (*[]dbsqlc.UserSession, error) {
	items, err := r.queries.DeleteExpiredUserSessions(ctx)
	if err != nil {
		return nil, err
	}
	return &items, nil
}