package middleware

import (
	"github.com/danivideda/satu-apotek-be/internal/repository"
	"github.com/danivideda/satu-apotek-be/internal/service"
)

type AppMiddleware struct {
	repo repository.Repository
	s    service.Session
}

func New(r repository.Repository, s service.Session) AppMiddleware {
	return AppMiddleware{
		repo: r,
		s:    s,
	}
}
