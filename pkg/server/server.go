// запуск сервака ListenAndServe
package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func Run(add string) error {
	log.Println("Запустился сервер")
	fmt.Println("Init server")
	r := chi.NewRouter()
	SetupRoutes()
	err := http.ListenAndServe(add, nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
		os.Exit(1)
		return err
	}

	return http.ListenAndServe(add, r)
}
