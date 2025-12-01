package handler

import (
	"github.com/danivideda/satu-apotek-be/internal/repository"
	"github.com/patrickmn/go-cache"
)

type Handler struct {
	Owner    *ownerHandler
	Auth     *authHandler
	Pharmacy *apotekHandler
	User     *userHandler
	AuthNew  *authNewHandler
}

func New(repo repository.Repository, cache *cache.Cache) Handler {
	return Handler{
		Auth:     &authHandler{repo: repo, cache: cache},
		Owner:    &ownerHandler{repo: repo, cache: cache},
		User:     &userHandler{repo: repo, cache: cache},
		Pharmacy: &apotekHandler{repo: repo, cache: cache},
		AuthNew:  &authNewHandler{repo: repo, cache: cache},
	}
}
