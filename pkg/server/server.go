package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Run(addr string) error {
	r := chi.NewRouter()
	SetupRoutes(r)

	log.Printf("Server starting on %s", addr)
	return http.ListenAndServe(addr, r)
}
