package handler

import (
	"github.com/danivideda/satu-apotek-be/internal/repository"
)

type Handler struct {
	Owner  *ownerHandler
	Auth   *authHandler
	User   *userHandler
	Pharmacy *pharmacyHandler
}

func New(repo repository.Repository) Handler {
	return Handler{
		Auth:   &authHandler{repo: repo},
		Owner:  &ownerHandler{repo: repo},
		User:   &userHandler{repo: repo},
		Pharmacy: &pharmacyHandler{repo: repo},
	}
}
