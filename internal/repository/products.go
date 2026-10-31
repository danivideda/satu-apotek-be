package repository

import (
	"context"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
)

type productsRepo struct {
	queries *dbsqlc.Queries
}

func (r *productsRepo) GetByPharmacyID(ctx context.Context, pharmacyID int64) (*[]dbsqlc.Product, error) {
	products, err := r.queries.GetProductsByPharmacyID(ctx, pharmacyID)
	if err != nil {
		if isNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	return &products, nil
}
