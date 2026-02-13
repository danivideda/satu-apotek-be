package middleware

import "github.com/danivideda/satu-apotek-be/internal/repository"

type AppMiddleware struct {
	repo repository.Repository
}

func New(r repository.Repository) AppMiddleware {
	return AppMiddleware{
		repo: r,
	}
}
