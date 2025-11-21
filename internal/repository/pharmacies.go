package repository

import (
	"context"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type PharmaciesRepo struct {
	queries *dbsqlc.Queries
}

func (r *PharmaciesRepo) Create(ctx context.Context, ownerID string, name string) (*dbsqlc.Pharmacy, error) {
	var ownerUUID pgtype.UUID
	if err := ownerUUID.Scan(ownerID); err != nil {
		return nil, err
	}
	params := dbsqlc.CreatePharmacyParams{
		OwnerID: ownerUUID,
		Name:    pgtype.Text{String: name, Valid: true},
	}
	pharmacy, err := r.queries.CreatePharmacy(ctx, params)
	if err != nil {
		return nil, err
	}

	return &pharmacy, nil
}
