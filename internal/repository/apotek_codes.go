package repository

import (
	"context"
	"time"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
	"github.com/danivideda/satu-apotek-be/internal/env"
	"github.com/jackc/pgx/v5/pgtype"
)

type pharmacyCodesRepo struct {
	queries *dbsqlc.Queries
}

func (r *pharmacyCodesRepo) Upsert(ctx context.Context, apotekID int64, code string) (*dbsqlc.PharmacyCode, error) {
	ttl, err := time.ParseDuration(env.GetString("CODE_TTL", "5m"))
	if err != nil {
		return nil, err
	}
	exp := pgtype.Timestamptz{
		Time:  time.Now().Add(ttl),
		Valid: true,
	}
	params := dbsqlc.UpsertApotekCodeParams{
		ApotekID:  apotekID,
		Code:      code,
		ExpiresAt: exp,
	}
	apotekCode, err := r.queries.UpsertApotekCode(ctx, params)
	if err != nil {
		return nil, err
	}

	return &apotekCode, nil
}

func (r *pharmacyCodesRepo) Get(ctx context.Context, apotekID int64) (*dbsqlc.PharmacyCode, error) {
	apotekCode, err := r.queries.GetApotekCode(ctx, apotekID)
	if err != nil {
		return nil, err
	}

	return &apotekCode, nil
}

func (r *pharmacyCodesRepo) GetByCode(ctx context.Context, code string) (*dbsqlc.PharmacyCode, error) {
	apotekCode, err := r.queries.GetApotekCodeByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	return &apotekCode, nil
}

func (r *pharmacyCodesRepo) DeleteExpired(ctx context.Context) (*[]dbsqlc.PharmacyCode, error) {
	apotekCode, err := r.queries.DeleteExpiredApotekCode(ctx)
	if err != nil {
		return nil, err
	}

	return &apotekCode, nil
}
