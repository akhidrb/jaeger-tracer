package main

import (
	middleware "github.com/MagalixTechnologies/core/middleware"
	"github.com/akhidrb/jaeger-tracer/mw"
	"github.com/go-chi/chi"
	"net/http"
)

func RunServer() {
	r := chi.NewRouter()
	r.Use(mw.SetServerSpan("server"))
	r.Use(middleware.Log(middleware.InfoLevel))

	r.Get("/publish", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	http.ListenAndServe("localhost:8082", r)
}

func main() {
	RunServer()
}
