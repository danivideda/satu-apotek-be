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
