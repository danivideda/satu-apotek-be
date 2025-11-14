package handler

import (
	"github.com/danivideda/satu-apotek-be/internal/store"
)

type Handler struct {
	Owner *ownerHandler
	Auth *authHandler
}

func New(store store.Storage) Handler {
	return Handler{
		Auth: &authHandler{store: store},
		Owner: &ownerHandler{store: store},
	}
}

type ownerHandler struct {
	store store.Storage
}

type authHandler struct {
	store store.Storage
}

