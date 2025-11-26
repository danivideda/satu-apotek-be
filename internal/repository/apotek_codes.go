package repository

import (
	"context"
	"time"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type apotekCodesRepo struct {
	queries *dbsqlc.Queries
}

func (r *apotekCodesRepo) Create(ctx context.Context, apotekID, code string) (*dbsqlc.ApotekCode, error) {
	var apotekUUID pgtype.UUID
	if err := apotekUUID.Scan(apotekID); err != nil {
		return nil, err
	}

	exp := pgtype.Timestamptz{
		Time:  time.Now().Add(1 * time.Minute),
		Valid: true,
	}
	params := dbsqlc.CreateApotekCodeParams{
		ApotekID: apotekUUID,
		Code:     code,
		Expires:  exp,
	}
	apotekCode, err := r.queries.CreateApotekCode(ctx, params)
	if err != nil {
		return nil, err
	}

	return &apotekCode, nil
}

func (r *apotekCodesRepo) Get(ctx context.Context, apotekID string) (*dbsqlc.ApotekCode, error) {
	var apotekUUID pgtype.UUID
	if err := apotekUUID.Scan(apotekID); err != nil {
		return nil, err
	}

	apotekCode, err := r.queries.GetApotekCode(ctx, apotekUUID)
	if err != nil {
		return nil, err
	}

	return &apotekCode, nil
}

func (r *apotekCodesRepo) GetByCode(ctx context.Context, code string) (*dbsqlc.ApotekCode, error) {
	apotekCode, err := r.queries.GetApotekCodeByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	return &apotekCode, nil
}
