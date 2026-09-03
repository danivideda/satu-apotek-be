package main

import (
	"fmt"
	"log"

	"github.com/danivideda/satu-apotek-be/internal/config"
	"github.com/danivideda/satu-apotek-be/internal/db"
	"github.com/danivideda/satu-apotek-be/internal/job"
	"github.com/danivideda/satu-apotek-be/internal/repository"
)

func main() {
	cfg := config.Load()

	fmt.Println(cfg.Job.Enabled)

	if !cfg.Job.Enabled {
		fmt.Println("Cron not enabled")
		return
	}

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

	select {}
}
