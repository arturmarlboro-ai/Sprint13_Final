package server

import (
	"Sprint13_Final/pkg/api"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r chi.Router) {
	// API
	r.Post("/api/task", api.PostTask)
	log.Println("Init routs")
	// Статика
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	})
	r.Handle("/js/*", http.StripPrefix("/js", http.FileServer(http.Dir("./web/js"))))
	r.Handle("/css/*", http.StripPrefix("/css", http.FileServer(http.Dir("./web/css"))))
	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/favicon.ico")
	})
}

//Сервер при запросе http://localhost:7540/ должен возвращать
// index.html из поддиректории web.
// Главная страница запрашивает .js и .css файлы,
// поэтому веб-сервер также должен их возвращать.
//
// Например:
// http://localhost:7540/js/scripts.min.js возвращает ./web/js/scripts.min.js;
// http://localhost:7540/css/style.css возвращает ./web/css/style.css;
// http://localhost:7540/favicon.ico возвращает ./web/favicon.ico.
