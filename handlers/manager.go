package handlers

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"miniapp/handlers/postgresql"
)

const (
	baseURL = "/miniapp"
)

func Manager(client *sql.DB) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route(baseURL, func(subRouter chi.Router) {
		repo := postgresql.NewPostgreSQLRepository(client)
		testHandler := NewHTTPHandler(repo)
		testHandler.RegisterRoutes(subRouter)
	})

	return r
}
