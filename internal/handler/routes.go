package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func CreateRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/trade", func(r chi.Router) {
		r.Get("/", GetTrade)
		r.Delete("/", DeleteAllTrade)
	})

	r.Route("/listener", func(r chi.Router) {
		r.Get("/{signer}", GetListener)
		r.Post("/", AddListener)
		r.Delete("/{signer}", RemoveListener)
	})

	return r
}
