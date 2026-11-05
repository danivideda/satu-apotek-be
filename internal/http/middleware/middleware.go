package middleware

import (
	"github.com/danivideda/satu-apotek-be/internal/config"
	c "github.com/danivideda/satu-apotek-be/internal/http/cookie"
	"github.com/danivideda/satu-apotek-be/internal/repository"
	"github.com/danivideda/satu-apotek-be/internal/service"
)

type AppMiddleware struct {
	repo       repository.Repository
	config     config.Config
	sessionSvc *service.Session
}

var cookie = c.New()

func New(r repository.Repository, cfg config.Config) AppMiddleware {
	sessionSvc := service.NewSession(cfg.Auth, r.CacheStore, r.OwnerSessions, r.UserSessions, r.PharmacySessions, r.Users, r.Pharmacies)
	return AppMiddleware{
		repo:       r,
		config:     cfg,
		sessionSvc: sessionSvc,
	}
}
