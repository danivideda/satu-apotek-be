package main

import (
	"log"

	"github.com/danivideda/satu-apotek-be/internal/db"
	"github.com/danivideda/satu-apotek-be/internal/env"
	"github.com/danivideda/satu-apotek-be/internal/http/handler"
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

	db, err := db.New(cfg.db.url)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	log.Println("Database connection pool established")

	r := repository.New(db)
	h := handler.New(r)
	s, err := job.NewScheduler(r)
	if err != nil {
		log.Panic(err)
	}
	s.AddClearCacheJob()
	s.Start()
	defer func() {
		if err := s.Shutdown(); err != nil {
			log.Panic(err)
		}
	}()

	app := &application{
		config:  cfg,
		handler: h,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
