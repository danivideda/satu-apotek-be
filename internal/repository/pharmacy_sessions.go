package repository

import (
	"context"
	"errors"
	"time"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type pharmacySessionsRepo struct {
	queries *dbsqlc.Queries
}

func (r *pharmacySessionsRepo) Create(ctx context.Context, pharmacyID int64, exp time.Time) (*dbsqlc.PharmacySession, error) {
	pharmacySession, err := r.queries.CreatePharmacySession(ctx, dbsqlc.CreatePharmacySessionParams{
		PharmacyID: pharmacyID,
		ExpiresAt:  pgtype.Timestamptz{Time: exp, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &pharmacySession, nil
}

func (r *pharmacySessionsRepo) Update(ctx context.Context, sessionID string, exp time.Time) (*dbsqlc.PharmacySession, error) {
	var sessionUUID pgtype.UUID
	if err := sessionUUID.Scan(sessionID); err != nil {
		return nil, err
	}

	pharmacySession, err := r.queries.UpdatePharmacySession(ctx, dbsqlc.UpdatePharmacySessionParams{
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

	return &pharmacySession, nil
}

func (r *pharmacySessionsRepo) Get(ctx context.Context, sessionID string) (*dbsqlc.PharmacySession, error) {
	var sessionUUID pgtype.UUID
	if err := sessionUUID.Scan(sessionID); err != nil {
		return nil, err
	}

	pharmacySession, err := r.queries.GetPharmacySession(ctx, sessionUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}
	return &pharmacySession, nil
}

func (r *pharmacySessionsRepo) Delete(ctx context.Context, sessionID string) (*dbsqlc.PharmacySession, error) {
	var sessionUUID pgtype.UUID
	if err := sessionUUID.Scan(sessionID); err != nil {
		return nil, err
	}
	deletedPharmacySession, err := r.queries.DeletePharmacySession(ctx, sessionUUID)
	if err != nil {
		return nil, err
	}
	return &deletedPharmacySession, nil
}

func (r *pharmacySessionsRepo) DeleteExpired(ctx context.Context) (*[]dbsqlc.PharmacySession, error) {
	items, err := r.queries.DeleteExpiredPharmacySessions(ctx)
	if err != nil {
		return nil, err
	}
	return &items, nil
}