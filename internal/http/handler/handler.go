package handler

import (
	"github.com/danivideda/satu-apotek-be/internal/repository"
)

type Handler struct {
	Owner    *ownerHandler
	Auth     *authHandler
	Pharmacy *apotekHandler
	User     *userHandler
}

func New(repo repository.Repository) Handler {
	return Handler{
		Auth:     &authHandler{repo: repo},
		Owner:    &ownerHandler{repo: repo},
		User:     &userHandler{repo: repo},
		Pharmacy: &apotekHandler{repo: repo},
	}
}
