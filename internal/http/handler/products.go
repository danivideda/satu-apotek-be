package handler

import (
	"net/http"

	"github.com/danivideda/satu-apotek-be/internal/http/json"
	"github.com/danivideda/satu-apotek-be/internal/http/middleware"
	"github.com/danivideda/satu-apotek-be/internal/repository"
)

type productHandler struct {
	repo repository.Repository
}

func newProductHandler(repo repository.Repository) *productHandler {
	return &productHandler{repo}
}

func (h *productHandler) GetByPharmacy(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authPharmacy, err := middleware.AuthPharmacyFromCtx(ctx)
	if err != nil {
		json.ResponseUnauthorized(w, r, err)
		return
	}

	products, err := h.repo.Products.GetByPharmacyID(ctx, authPharmacy.PharmacyID)
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
