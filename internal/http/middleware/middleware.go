package middleware

import (
	"github.com/danivideda/satu-apotek-be/internal/config"
	c "github.com/danivideda/satu-apotek-be/internal/http/cookie"
	"github.com/danivideda/satu-apotek-be/internal/repository"
)

type AppMiddleware struct {
	repo   repository.Repository
	config config.Config
}

var cookie = c.New()

func New(r repository.Repository, cfg config.Config) AppMiddleware {
	return AppMiddleware{
		repo:   r,
		config: cfg,
	}
}
