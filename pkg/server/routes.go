package server

import (
	"fmt"
	"log"
	"net/http"
)

const (
	webDir1 = "./web"
	webDir2 = "./web/js/scripts.min.js"
	webDir3 = "./web/css/style.css"
	webDir4 = "./web"
)

func SetupRoutes() {
	log.Println("Init routes")
	fmt.Println("Init routs")
	http.Handle("/", http.FileServer(http.Dir(webDir1)))
	http.Handle("/js", http.FileServer(http.Dir(webDir2)))
	http.Handle("/css", http.FileServer(http.Dir(webDir3)))
	http.Handle("/favicon.ico", http.FileServer(http.Dir(webDir4)))
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
