package handler

import (
	"github.com/danivideda/satu-apotek-be/internal/store"
)

type Handler struct {
	Owner *ownerHandler
	Auth  *authHandler
	Pharmacy *apotekHandler
}

func New(store store.Storage) Handler {
	return Handler{
		Auth:  &authHandler{store: store},
		Owner: &ownerHandler{store: store},
		Pharmacy: &apotekHandler{store: store},
	}
}
