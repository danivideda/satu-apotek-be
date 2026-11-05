package main

import (
	"log"

	"github.com/danivideda/satu-apotek-be/internal/config"
	"github.com/danivideda/satu-apotek-be/internal/db"
	"github.com/danivideda/satu-apotek-be/internal/http/handler"
	"github.com/danivideda/satu-apotek-be/internal/http/middleware"
	"github.com/danivideda/satu-apotek-be/internal/repository"
	"github.com/danivideda/satu-apotek-be/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := db.NewPostgres(cfg.DB.URL)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	log.Println("Database connection pool established")

	c, err := repository.NewCacheStore()
	if err != nil {
		log.Panic(err)
	}

	r := repository.New(db, c)

	sessionSvc := service.NewSession(cfg.Auth, r.CacheStore, r.OwnerSessions, r.UserSessions, r.PharmacySessions, r.Users, r.Pharmacies)
	authSvc := service.NewAuth(r.Owners, r.Users, r.Pharmacies, sessionSvc)
	productSvc := service.NewProduct(r.Products)

	h := handler.New(r, cfg, authSvc, productSvc)
	md := middleware.New(r, cfg)

	app := &application{
		config:     cfg,
		handler:    h,
		middleware: md,
	}

	mux := app.mount(cfg.CORS)
	log.Fatal(app.run(mux))
}
