package server

import (
	"net/http"
)

func GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Task ID: " + id))
}
func SecondHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("This is the second handler."))
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
