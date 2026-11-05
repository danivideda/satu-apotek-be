package service

import (
	"context"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
	"github.com/danivideda/satu-apotek-be/internal/repository"
)

type Product struct {
	productsRepo repository.ProductsRepository
}

func NewProduct(productsRepo repository.ProductsRepository) *Product {
	return &Product{productsRepo}
}

func (p *Product) GetAllByPharmacyID(ctx context.Context, pharmacyID int64) (*[]dbsqlc.Product, error) {
	products, err := p.productsRepo.GetByPharmacyID(ctx, pharmacyID)
	return products, err
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

func (p *Product) Create(ctx context.Context, in CreateProductInput) (*dbsqlc.Product, error) {
	product, err := p.productsRepo.Create(ctx, repository.CreateProductInput{
		PharmacyID:    in.PharmacyID,
		Name:          in.Name,
		SKU:           in.SKU,
		Barcode:       in.Barcode,
		CategoryID:    in.CategoryID,
		BaseUnitID:    in.BaseUnitID,
		MinStock:      in.MinStock,
		BaseSellPrice: in.BaseSellPrice,
		IsActive:      in.IsActive,
	})

	return product, err
}
