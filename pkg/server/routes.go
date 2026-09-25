package server

import (
	"Sprint13_Final/pkg/api"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r chi.Router) {
	r.Post("/api/task", api.PostTask)
	log.Println("Init routs")

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	})
	r.Handle("/*", http.FileServer(http.Dir("./web")))

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
