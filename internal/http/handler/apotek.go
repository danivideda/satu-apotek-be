package handler

import (
	"net/http"

	"github.com/danivideda/satu-apotek-be/internal/http/json"
	"github.com/danivideda/satu-apotek-be/internal/http/middleware"
	"github.com/danivideda/satu-apotek-be/internal/http/response"
	"github.com/danivideda/satu-apotek-be/internal/store"
)

type apotekHandler struct {
	store store.Storage
}

func (h *apotekHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	type createPharmacyPayload struct {
		Name string `json:"name"`
	}
	var payload createPharmacyPayload
	if err := json.Read(w, r, &payload); err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	auth := middleware.AuthClaimsFromContext(ctx)
	if auth == nil {
		response.InternalServerError(w, r, ErrInvalidAuthToken)
		return 
	}

	pharmacy, err := h.store.Pharmacies.Create(ctx, auth.ID, payload.Name)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	if err := response.Created(w, r, map[string]any{"pharmacy": pharmacy}); err != nil {
		response.InternalServerError(w, r, err)
		return
	}
}