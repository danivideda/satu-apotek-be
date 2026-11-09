package repository

import (
	"context"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
)

type productUnitsRepo struct {
	queries *dbsqlc.Queries
}

func (r *productUnitsRepo) GetByPharmacyID(ctx context.Context, pharmacyID int64) (*[]dbsqlc.ProductUnit, error) {
	productUnits, err := r.queries.ProductUnitsGetByPharmacyID(ctx, pharmacyID)
	if err != nil {
		if isNotFoundError(err) {
			emptyArr := make([]dbsqlc.ProductUnit, 0)
			return &emptyArr, nil
		}
		return nil, err
	}

	return &productUnits, nil
}

func (r *productUnitsRepo) Create(ctx context.Context, pharmacyID int64, name, label string) (*dbsqlc.ProductUnit, error) {
	productUnit, err := r.queries.ProductUnitsCreate(ctx, dbsqlc.ProductUnitsCreateParams{
		PharmacyID: pharmacyID,
		Name:       name,
		Label:      label,
	})
	if err != nil {
		if isDuplicateError(err) {
			return nil, ErrDuplicateValue
		}
		return nil, err
	}

	return &productUnit, nil
}
