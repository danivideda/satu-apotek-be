package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/danivideda/satu-apotek-be/internal/http/json"
	"github.com/danivideda/satu-apotek-be/internal/http/middleware"
	"github.com/danivideda/satu-apotek-be/internal/repository"
	"github.com/danivideda/satu-apotek-be/internal/service"
	"github.com/go-chi/chi/v5"
)

type productHandler struct {
	productService *service.Product
}

const productCtx = "ProductCtx"

func newProductHandler(productService *service.Product) *productHandler {
	return &productHandler{productService}
}

func (h *productHandler) ProductCtx(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// 1. Get productID
		productIDString := chi.URLParam(r, "productID")
		if productIDString == "" {
			json.ResponseBadRequest(w, r, ErrProductWrongFormat)
			return
		}

		productID, err := strconv.ParseInt(productIDString, 10, 64)
		if err != nil {
			json.ResponseBadRequest(w, r, ErrProductWrongFormat)
			return
		}

		// 2. Get product details
		product, err := h.productService.GetDetailByID(ctx, productID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				json.ResponseNotFound(w, r, err)
				return
			}
			json.ResponseInternalServerError(w, r, err)
			return
		}

		// 3. Set the context
		ctx = context.WithValue(ctx, productCtx, product)

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

func (h *productHandler) ProductFromCtx(ctx context.Context) (*service.ProductDetail, error) {
	product, ok := ctx.Value(productCtx).(*service.ProductDetail)
	if !ok {
		return nil, errors.New("ProductCtx type assertion missmatch")
	}
	return product, nil
}

func (h *productHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	productDetail, err := h.ProductFromCtx(ctx)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
	if err := json.ResponseOK(w, productDetail); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}

func (h *productHandler) GetAllByPharmacyID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authPharmacy, err := middleware.AuthPharmacyFromCtx(ctx)
	if err != nil {
		json.ResponseUnauthorized(w, r, err)
		return
	}

	// Get URL query for pagination
	query := r.URL.Query()
	var offset int32 = 0
	var limit int32 = 10
	if query.Has("offset") && query.Has("limit") {
		offsetInt, err := strconv.Atoi(query.Get("offset"))
		if err != nil {
			json.ResponseInternalServerError(w, r, err)
			return
		}
		offset = int32(offsetInt)

		limitInt, err := strconv.Atoi(query.Get("limit"))
		if err != nil {
			json.ResponseInternalServerError(w, r, err)
			return
		}
		limit = int32(limitInt)
	}

	products, err := h.productService.GetAllByPharmacyID(ctx, authPharmacy.PharmacyID, offset, limit)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			json.ResponseNotFound(w, r, err)
			return
		}
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
		CategoryID    *int64  `json:"category_id" validate:"required"`
		BaseUnitID    *int64  `json:"base_unit_id" validate:"required"`
		MinStock      int32   `json:"min_stock"`
		BaseSellPrice *string `json:"base_sell_price"`
		IsActive      bool    `json:"is_active"`
		// Units         []any   `json:"unit_conversions"`
	}
	if ok := parseAndValidateJSONPayload(w, r, &payload); !ok {
		return
	}

	newProduct, err := h.productService.Create(ctx, service.CreateProductInput{
		PharmacyID:    authPharmacy.PharmacyID,
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

func (h *productHandler) GetAllUnitByPharmacyID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authPharmacy, err := middleware.AuthPharmacyFromCtx(ctx)
	if err != nil {
		json.ResponseUnauthorized(w, r, err)
		return
	}

	productUnits, err := h.productService.GetAllUnitByPharmacyID(ctx, authPharmacy.PharmacyID)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	if err := json.ResponseOK(w, productUnits); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}

func (h *productHandler) CreateUnit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authPharmacy, err := middleware.AuthPharmacyFromCtx(ctx)
	if err != nil {
		json.ResponseUnauthorized(w, r, err)
		return
	}

	var payload struct {
		Name  string `json:"name" validate:"required"`
		Label string `json:"label" validate:"required"`
	}
	if ok := parseAndValidateJSONPayload(w, r, &payload); !ok {
		return
	}
	newUnit, err := h.productService.CreateUnitByPharmacyID(ctx, authPharmacy.PharmacyID, payload.Name, payload.Label)
	if err != nil {
		if errors.Is(err, service.ErrCreateItem) {
			json.ResponseBadRequest(w, r, err)
			return
		}
		json.ResponseInternalServerError(w, r, err)
		return
	}

	if err := json.ResponseCreated(w, newUnit); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}

func (h *productHandler) CreateUnitConversion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	productDetail, err := h.ProductFromCtx(ctx)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	var payload struct {
		UnitID             int64 `json:"unit_id" validate:"required"`
		QuantityInBaseUnit int32 `json:"quantity_in_base_unit" validate:"required"`
	}
	if ok := parseAndValidateJSONPayload(w, r, &payload); !ok {
		return
	}

	productUnitConversion, err := h.productService.CreateConversionForProduct(ctx, service.CreateConversionForProductInput{
		ProductID:          productDetail.ID,
		UnitID:             payload.UnitID,
		QuantityInBaseUnit: payload.QuantityInBaseUnit,
	})

	if err != nil {
		if errors.Is(err, service.ErrCreateItem) {
			json.ResponseBadRequest(w, r, err)
			return
		}
		json.ResponseInternalServerError(w, r, err)
		return
	}

	if err := json.ResponseCreated(w, productUnitConversion); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}
