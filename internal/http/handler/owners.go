package handler

import (
	"log"
	"net/http"

	"github.com/danivideda/satu-apotek-be/internal/http/json"
)


func (h *Handler) GetOwnerByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ownerID := 3
	owners, err := h.store.Owners.GetOwnerByID(ctx, ownerID)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	if err := json.ResponseJSON(w, http.StatusCreated, owners); err != nil {
		log.Printf("Error occured: %s", err)
	}
}
