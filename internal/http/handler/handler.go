package handler

import (
	"github.com/danivideda/satu-apotek-be/internal/repository"
)

type Handler struct {
	Owner  *ownerHandler
	Auth   *authHandler
	User   *userHandler
	Apotek *apotekHandler
}

func New(repo repository.Repository) Handler {
	return Handler{
		Auth:   &authHandler{repo: repo},
		Owner:  &ownerHandler{repo: repo},
		User:   &userHandler{repo: repo},
		Apotek: &apotekHandler{repo: repo},
	}
}
