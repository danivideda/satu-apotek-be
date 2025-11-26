package repository

import (
	"context"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type pharmaciesRepo struct {
	queries *dbsqlc.Queries
}

func (r *pharmaciesRepo) Create(ctx context.Context, ownerID string, name string) (*dbsqlc.Pharmacy, error) {
	var ownerUUID pgtype.UUID
	if err := ownerUUID.Scan(ownerID); err != nil {
		return nil, err
	}
	params := dbsqlc.CreatePharmacyParams{
		OwnerID: ownerUUID,
		Name:    name,
	}
	pharmacy, err := r.queries.CreatePharmacy(ctx, params)
	if err != nil {
		return nil, err
	}

	return &pharmacy, nil
}
