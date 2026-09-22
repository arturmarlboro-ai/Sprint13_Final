package api

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
