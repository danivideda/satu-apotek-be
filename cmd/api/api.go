package main

import (
	"log"
	"net/http"
	"time"

	"github.com/danivideda/satu-apotek-be/internal/http/handler"
	"github.com/danivideda/satu-apotek-be/internal/http/middleware"
	"github.com/go-chi/chi/v5"
	cm "github.com/go-chi/chi/v5/middleware"
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
	r.Use(cm.RequestID)
	r.Use(cm.RealIP)
	r.Use(cm.Logger)
	r.Use(cm.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(cm.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world!"))
	})

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Health check: OK\n"))
		})

		r.Route("/auth", func(r chi.Router) {
			r.Route("/owners", func(r chi.Router) {
				r.Post("/register", app.handler.Auth.OwnerRegister)
				r.Post("/login", app.handler.Auth.OwnerLogin)
				r.Get("/refresh", app.handler.Auth.OwnerRefresh)
				r.With(middleware.AuthOwner).Post("/logout", app.handler.Auth.OwnerLogout)
			})
			r.Route("/users", func(r chi.Router) {
				r.Post("/login", app.handler.Auth.UserLogin)
				r.Get("/refresh", app.handler.Auth.UserRefresh)
				r.With(middleware.AuthOwner).Post("/logout", app.handler.Auth.UserLogout)
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(middleware.AuthOwner)
				r.Post("/create", app.handler.User.Create)
			})
		})

		r.Route("/owners", func(r chi.Router) {
			r.Get("/", app.handler.Owner.GetByID)
		})

		r.Route("/apotek", func(r chi.Router) {
			r.Route("/code", func(r chi.Router) {
				r.Post("/verify", app.handler.Pharmacy.VerifyCode)
				r.With(middleware.AuthOwner).Post("/create", app.handler.Pharmacy.CreateCode)
			})

			r.Group(func(r chi.Router) {
				r.Use(middleware.AuthOwner)
				r.Post("/create", app.handler.Pharmacy.Create)
			})
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
