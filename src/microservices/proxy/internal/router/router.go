package router

import (
	"proxy/internal/application"
	"proxy/internal/handler"

	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
)

func New(app *application.App) *chi.Mux {
	r := chi.NewRouter()
	r.Use(chi_middleware.Logger)

	app.Log.Info("Server started")

	monolithHandler := handler.NewProxyHandler(app, app.Conf.MonolithUrl)
	r.HandleFunc("/*", monolithHandler.Proxy)

	r.Mount("/api/movies", movies(app))

	return r
}
