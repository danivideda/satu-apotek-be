package service

import (
	"context"
	"errors"
	"fmt"
	"time"

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

type ProductDetail struct {
	*dbsqlc.ProductsGetByIDRow
}

func (p *Product) GetDetailByID(ctx context.Context, productID int64) (*ProductDetail, error) {
	product, err := p.productsRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	productDetail := ProductDetail{product}
	return &productDetail, nil
}

type ProductItem struct {
	ID            int64         `json:"id"`
	Name          string        `json:"name"`
	Sku           *string       `json:"sku"`
	Barcode       *string       `json:"barcode"`
	MinStock      int32         `json:"min_stock"`
	BaseSellPrice string        `json:"base_sell_price"`
	IsActive      bool          `json:"is_active"`
	Category      *CategoryItem `json:"category"`
	BaseUnit      *BaseUnitItem `json:"base_unit"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

type CategoryItem struct {
	ID   int64  `json:"id,omitempty"`
	Code string `json:"code,omitempty"`
	Name string `json:"name,omitempty"`
}
type BaseUnitItem struct {
	ID    int64  `json:"id,omitempty"`
	Label string `json:"label,omitempty"`
	Name  string `json:"name,omitempty"`
}

func (p *Product) GetAllByPharmacyID(ctx context.Context, pharmacyID int64, offset, limit int32) (*[]ProductItem, error) {
	products, err := p.productsRepo.GetByPharmacyID(ctx, pharmacyID, offset, limit)
	if err != nil {
		return nil, err
	}

	var productList []ProductItem
	for _, item := range *products {
		productItem := ProductItem{
			ID:   item.Product.ID,
			Name: item.Product.Name,
			// Sku:           &item.CategoryCode.String,
			// Barcode:       &item.Product.Barcode.String,
			MinStock:      item.Product.MinStock,
			BaseSellPrice: numericToString(item.Product.BaseSellPrice),
			IsActive:      item.Product.IsActive,
			// Category:      &CategoryItem{},
			// BaseUnit:      &BaseUnitItem{},
			CreatedAt: item.Product.CreatedAt.Time,
			UpdatedAt: item.Product.UpdatedAt.Time,
		}

		if item.Product.Sku.Valid {
			productItem.Sku = &item.Product.Sku.String
		}

		if item.Product.Barcode.Valid {
			productItem.Barcode = &item.Product.Barcode.String
		}

		if item.Product.CategoryID.Valid {
			productItem.Category = &CategoryItem{
				ID:   item.Product.CategoryID.Int64,
				Code: item.CategoryCode.String,
				Name: item.CategoryName.String,
			}
		}

		if item.Product.BaseUnitID.Valid {
			productItem.BaseUnit = &BaseUnitItem{
				ID:    item.Product.BaseUnitID.Int64,
				Label: item.UnitLabel.String,
				Name:  item.UnitName.String,
			}
		}

		productList = append(productList, productItem)
	}

	return &productList, err
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

type CreateConversionForProductInput struct {
	ProductID          int64 `json:"product_id" validate:"required,number"`
	UnitID             int64 `json:"unit_id" validate:"required,number"`
	QuantityInBaseUnit int32 `json:"quantity_in_base_unit" validate:"required,number"`
}

func (p *Product) CreateConversionForProduct(ctx context.Context, input CreateConversionForProductInput) (*dbsqlc.ProductUnitConversion, error) {
	productUnitConversions, err := p.productUnitsRepo.CreateConversionForProduct(ctx, input.ProductID, input.UnitID, input.QuantityInBaseUnit)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateValue) {
			return nil, fmt.Errorf("%w: %w", ErrCreateItem, err)
		}
		return nil, err
	}

	return productUnitConversions, nil
}
