package handler

import (
	"net/http"

	"github.com/danivideda/satu-apotek-be/internal/http/json"
	"github.com/danivideda/satu-apotek-be/internal/http/response"
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
		response.BadRequest(w, r, err)
		return
	}

	if err := json.WriteResponse(w, http.StatusOK, owners); err != nil {
		response.InternalServerError(w, r, err)
		return
	}
}
