package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {

	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-type", "X-CSRF-Token"},
		AllowCredentials: true,
		ExposedHeaders: []string{"Link"},
		MaxAge: 300,
	}))

	mux.Use(middleware.Heartbeat("/ping"));

	mux.Post("/", app.Broker)
	mux.Post("/handle", app.HandleSubmission)
	mux.Post("/log-grpc", app.LogVia_gRPC)

	return mux;
}