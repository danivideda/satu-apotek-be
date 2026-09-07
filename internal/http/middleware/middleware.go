package middleware

import (
	"github.com/danivideda/satu-apotek-be/internal/config"
	"github.com/danivideda/satu-apotek-be/internal/repository"
	"github.com/danivideda/satu-apotek-be/internal/service"
)

type AppMiddleware struct {
	repo   repository.Repository
	s      service.Session
	config config.Config
}

func New(r repository.Repository, s service.Session, cfg config.Config) AppMiddleware {
	return AppMiddleware{
		repo:   r,
		s:      s,
		config: cfg,
	}
}
