package main

import (
	"log"
	"net/http"
	"time"

	"github.com/danivideda/satu-apotek-be/internal/http/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config  config
	handler handler.Handler
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	url string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world!"))
	})

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Health check: OK\n"))
		})

		r.Route("/owners", func(r chi.Router) {
			r.Get("/", app.handler.GetOwnerByID)
			r.Post("/create", app.handler.CreateOwner)
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	srv := http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Listening on http://%s", app.config.addr)

	return srv.ListenAndServe()
}
