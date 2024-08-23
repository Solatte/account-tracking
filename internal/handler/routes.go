package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func CreateRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	var TradeHandler = NewTradeHandler()
	var ListenerHandler = NewListenerHandler()

	ListenerHandler.Init()

	r.Route("/trade", func(r chi.Router) {
		r.Get("/", TradeHandler.Get)
		r.Delete("/", TradeHandler.DeleteAll)
	})

	r.Route("/listener", func(r chi.Router) {
		r.Get("/{signer}", ListenerHandler.Get)
		r.Post("/", ListenerHandler.Add)
		r.Delete("/{signer}", ListenerHandler.Remove)
	})

	return r
}
