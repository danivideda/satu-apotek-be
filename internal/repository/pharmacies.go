package repository

import (
	"context"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
)

type pharmaciesRepo struct {
	queries *dbsqlc.Queries
}

func (r *pharmaciesRepo) Create(ctx context.Context, ownerID int64, name string) (*dbsqlc.Pharmacy, error) {
	params := dbsqlc.CreatePharmacyParams{OwnerID: ownerID, Name: name}
	pharmacy, err := r.queries.CreatePharmacy(ctx, params)
	if err != nil {
		return nil, err
	}

	return &pharmacy, nil
}

func (r *pharmaciesRepo) GetByOwner(ctx context.Context, ownerID int64) (*[]dbsqlc.Pharmacy, error) {
	pharmacies, err := r.queries.GetPharmaciesByOwner(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	return &pharmacies, nil
}