package handler

import "github.com/danivideda/satu-apotek-be/internal/repository"

type productUnitHandler struct {
	productUnitsRepo repository.ProductUnitsRepository
}

func newProductUnitHandler(productUnitsRepo repository.ProductUnitsRepository) *productUnitHandler {
	return &productUnitHandler{productUnitsRepo}
}

