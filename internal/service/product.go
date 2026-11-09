package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
	"github.com/danivideda/satu-apotek-be/internal/repository"
)

type Product struct {
	productsRepo     repository.ProductsRepository
	productUnitsRepo repository.ProductUnitsRepository
}

func NewProduct(productsRepo repository.ProductsRepository, productUnitsRepo repository.ProductUnitsRepository) *Product {
	return &Product{productsRepo, productUnitsRepo}
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

func (p *Product) GetAllUnitByPharmacyID(ctx context.Context, pharmacyID int64) (*[]dbsqlc.ProductUnit, error) {
	productUnits, err := p.productUnitsRepo.GetByPharmacyID(ctx, pharmacyID)
	return productUnits, err
}

func (p *Product) CreateUnitByPharmacyID(ctx context.Context, pharmacyID int64, name, label string) (*dbsqlc.ProductUnit, error) {
	productUnit, err := p.productUnitsRepo.Create(ctx, pharmacyID, name, label)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateValue) {
			return nil, fmt.Errorf("%w: %w", ErrCreateItem, err)
		}
		return nil, err
	}
	return productUnit, nil
}
