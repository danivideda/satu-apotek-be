package main

import (
	"log"

	"github.com/danivideda/satu-apotek-be/internal/db"
	"github.com/danivideda/satu-apotek-be/internal/env"
	"github.com/danivideda/satu-apotek-be/internal/http/handler"
	"github.com/danivideda/satu-apotek-be/internal/http/middleware"
	"github.com/danivideda/satu-apotek-be/internal/job"
	"github.com/danivideda/satu-apotek-be/internal/repository"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", "localhost:8080"),
		db: dbConfig{
			url: env.GetString("DATABASE_URL", "admin:adminpassword@localhost/satuapotek?sslmode=disable"),
		},
	}

	db, err := db.NewPostgres(cfg.db.url)
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
	h := handler.New(r)
	md := middleware.New(r)

	// TODO: Move this to separate process /cmd/cron/main.go
	// Cron Job
	s, err := job.NewScheduler(r)
	if err != nil {
		log.Panic(err)
	}
	s.AddClearApotekCodeJob()
	s.AddDeleteExpiredSessionsJob()
	s.Start()
	defer func() {
		if err := s.Shutdown(); err != nil {
			log.Panic(err)
		}
	}()


	app := &application{
		config:     cfg,
		handler:    h,
		middleware: md,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
