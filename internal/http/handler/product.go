package handler

import (
	"net/http"

	"github.com/danivideda/satu-apotek-be/internal/http/json"
	"github.com/danivideda/satu-apotek-be/internal/http/middleware"
	"github.com/danivideda/satu-apotek-be/internal/service"
)

type productHandler struct {
	productService *service.Product
}

func newProductHandler(productService *service.Product) *productHandler {
	return &productHandler{productService}
}

func (h *productHandler) GetAllByPharmacyID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authPharmacy, err := middleware.AuthPharmacyFromCtx(ctx)
	if err != nil {
		json.ResponseUnauthorized(w, r, err)
		return
	}

	products, err := h.productService.GetAllByPharmacyID(ctx, authPharmacy.PharmacyID)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	res := products
	if err := json.ResponseOK(w, res); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}

func (h *productHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authPharmacy, err := middleware.AuthPharmacyFromCtx(ctx)
	if err != nil {
		json.ResponseUnauthorized(w, r, err)
		return
	}

	// Get payload of product information
	var payload struct {
		Name          string  `json:"name" validate:"required"`
		SKU           *string `json:"sku"`
		Barcode       *string `json:"barcode"`
		CategoryID    *int64  `json:"category_id"`
		BaseUnitID    *int64  `json:"base_unit_id"`
		MinStock      int32   `json:"min_stock"`
		BaseSellPrice *string `json:"base_sell_price"`
		IsActive      bool    `json:"is_active"`
	}
	if ok := parseAndValidateJSONPayload(w, r, &payload); !ok {
		return
	}

	newProduct, err := h.productService.Create(ctx, service.CreateProductInput{
		PharmacyID: authPharmacy.PharmacyID, 
		Name:          payload.Name,
		SKU:           payload.SKU,
		Barcode:       payload.Barcode,
		CategoryID:    payload.CategoryID,
		BaseUnitID:    payload.BaseUnitID,
		MinStock:      payload.MinStock,
		BaseSellPrice: payload.BaseSellPrice,
		IsActive:      payload.IsActive,
	})
	if err != nil {
		json.ResponseBadRequest(w, r, err)
		return
	}

	res := newProduct
	if err := json.ResponseOK(w, res); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}
