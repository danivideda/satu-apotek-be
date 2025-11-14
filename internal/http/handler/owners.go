package handler

import (
	"net/http"

	"github.com/danivideda/satu-apotek-be/internal/http/json"
)

func (h *ownerHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ownerID := 3
	owners, err := h.store.Owners.GetOwnerByID(ctx, ownerID)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	if err := json.WriteResponse(w, http.StatusOK, owners); err != nil {
		internalServerErrorResponse(w, r, err)
		return
	}
}
