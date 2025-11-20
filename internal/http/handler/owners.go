package handler

import (
	"net/http"

	"github.com/danivideda/satu-apotek-be/internal/http/json"
	"github.com/danivideda/satu-apotek-be/internal/store"
)

type ownerHandler struct {
	store store.Storage
}

func (h *ownerHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ownerID := "todo-uuid"
	owners, err := h.store.Owners.GetByID(ctx, ownerID)
	if err != nil {
		json.ResponseBadRequest(w, r, err)
		return
	}

	if err := json.WriteResponse(w, http.StatusOK, owners); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}
