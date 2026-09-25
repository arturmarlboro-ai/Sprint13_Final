package api

import (
	"encoding/json"
	"log"
	"net/http"

	"Sprint13_Final/pkg/db"
)

func PostTask(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Printf("PostTask error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		log.Println("task.Title emtpy")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(r.Context(), task)
	if err != nil {
		log.Printf("AddTask error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}
