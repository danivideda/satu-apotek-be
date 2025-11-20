package handler

import (
	"github.com/danivideda/satu-apotek-be/internal/store"
)

type Handler struct {
	Owner *ownerHandler
	Auth  *authHandler
	Pharmacy *apotekHandler
	User *userHandler
}

func New(store store.Storage) Handler {
	return Handler{
		Auth:  &authHandler{store: store},
		Owner: &ownerHandler{store: store},
		User: &userHandler{store: store},
		Pharmacy: &apotekHandler{store: store},
	}
}
