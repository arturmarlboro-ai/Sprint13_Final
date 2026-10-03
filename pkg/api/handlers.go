package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"Sprint13_Final/pkg/db"
)

func PostTask(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	now := time.Now()

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

	dstart := task.Date
	repeat := task.Repeat

	Nextdate, err := NextDate(now, dstart, repeat)
	if err != nil {
		log.Printf("error NextDate: %v", err)
	}
	fmt.Println(Nextdate)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}
