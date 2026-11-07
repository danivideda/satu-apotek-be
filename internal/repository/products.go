package repository

import (
	"context"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type productsRepo struct {
	queries *dbsqlc.Queries
}

func (r *productsRepo) GetByPharmacyID(ctx context.Context, pharmacyID int64) (*[]dbsqlc.Product, error) {
	products, err := r.queries.ProductsGetByPharmacyID(ctx, pharmacyID)
	if err != nil {
		if isNotFoundError(err) {
			emptyArr := make([]dbsqlc.Product, 0)
			return &emptyArr, nil
		}
		return nil, err
	}

	return &products, nil
}

type CreateProductInput struct {
	PharmacyID    int64
	Name          string
	SKU           *string
	Barcode       *string
	CategoryID    *int64
	BaseUnitID    *int64
	MinStock      int32
	BaseSellPrice *string
	IsActive      bool
}

func (r *productsRepo) Create(ctx context.Context, in CreateProductInput) (*dbsqlc.Product, error) {
	params := dbsqlc.ProductsCreateParams{
		PharmacyID: in.PharmacyID,
		Name:       in.Name,
		MinStock:   in.MinStock,
		IsActive:   in.IsActive,
	}

	if in.SKU != nil {
		params.Sku = pgtype.Text{String: *in.SKU, Valid: true}
	}
	if in.Barcode != nil {
		params.Barcode = pgtype.Text{String: *in.Barcode, Valid: true}
	}
	if in.CategoryID != nil {
		params.CategoryID = pgtype.Int8{Int64: *in.CategoryID, Valid: true}
	}
	if in.BaseUnitID != nil {
		params.BaseUnitID = pgtype.Int8{Int64: *in.BaseUnitID, Valid: true}
	}
	if in.BaseSellPrice != nil {
		var n pgtype.Numeric
		if err := n.Scan(*in.BaseSellPrice); err != nil {
			return nil, err
		}
		params.BaseSellPrice = n
	}

	product, err := r.queries.ProductsCreate(ctx, params)
	if err != nil {
		return nil, err
	}

	return &product, nil
}
