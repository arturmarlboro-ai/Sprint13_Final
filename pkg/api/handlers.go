package api

import (
	"Sprint13_Final/pkg/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func nextDayHandler(w http.ResponseWriter, r *http.Request) {

	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeatStr := r.URL.Query().Get("repeat")

	var now time.Time

	if nowStr == "" {
		now = time.Now()
	}

	now, err := time.Parse(DateFormat, nowStr)
	if err != nil {
		log.Printf("error time parse at handlers.go")
		http.Error(w, fmt.Sprintf("invalid 'date' format, use %s", DateFormat), http.StatusBadRequest)
		return
	}

	nextDateStr, err := NextDate(now, dateStr, repeatStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "%s", nextDateStr)
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	var tasks []*db.Task
	var err error
	search := r.URL.Query().Get("search")
	//search = strings.ToLower(search)

	fmt.Println(search)
	log.Print(search)

	if search == "" {
		tasks, err = db.GetTasksByLimit(r.Context(), 50)
		if err != nil {
			writeJson(w, ("error: failed to get tasks"))
			return
		}

		if tasks == nil {

			writeJson(w, TasksResp{
				Tasks: tasks,
			})
			return
		}
	}

	tasks, err = db.GetTasksBySearch(r.Context(), search)
	if err != nil {
		writeJson(w, ("error: failed to get tasks"))
		return
	}

	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Printf("PostTask json decode error: %s", err)
		writeJson(w, "error: invalid JSON")
		return
	}

	if task.Title == "" {
		log.Println("task.Title emtpy")
		writeJson(w, "error: title is required")
		return
	}

	if err := checkDate(&task); err != nil {
		log.Printf("error checkDate: %v", err)
		writeJson(w, map[string]string{"error": fmt.Sprintf("invalid date: %v", err)})
		return
	}

	id, err := db.AddTask(r.Context(), task)
	if err != nil {
		log.Printf("AddTask error: %v", err)
		writeJson(w, map[string]string{"error": "failed to add task"})
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//w.WriteHeader(http.StatusCreated)
	//writeJson(w, fmt.Sprintf("id: %s", id))
	writeJson(w, map[string]string{"id": id})

}

func editTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	task, err := db.GetTaskByID(r.Context(), id)
	if err != nil {
		log.Printf("editTaskHandler error: %v", err)
		http.NotFound(w, r)
		writeJson(w, map[string]string{"error": "failed get task by id"})
		return
	}
	log.Print(task)
	writeJson(w, task)
}

func putTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Printf("putTask json decode error: %v", err)
		writeJson(w, map[string]string{"error": "invalid json"})
		return
	}
	err := db.UpdateTask(r.Context(), task, task.ID)
	if err != nil {
		log.Printf("putTask erroe: %v", err)
		writeJson(w, map[string]string{"error": "failed to update task"})
		return
	}

	writeJson(w, map[string]string{})
}

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	task, err := db.GetTaskByID(r.Context(), id)
	if err != nil {
		log.Printf("doneTaskHandler error: %v", err)
		http.NotFound(w, r)
		writeJson(w, map[string]string{"error": "failed get task by id"})
		return
	}
	if task.Repeat != "" {
		nextDateStr, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			log.Printf("doneTaskHandler error: %v", err)
			writeJson(w, map[string]string{"error": "failed to calculate next date"})
			return
		}
		task.Date = nextDateStr

		err = db.UpdateTask(r.Context(), *task, task.ID)
		if err != nil {
			log.Printf("doneTaskHandler error: %v", err)
			writeJson(w, map[string]string{"error": "failed to update task with next date"})
			return
		}
	}

	err = db.DeleteTask(r.Context(), id)
	if err != nil {
		log.Printf("doneTaskHandler error: %v", err)
		writeJson(w, map[string]string{"error": "failed to delete task"})
		return
	}
	writeJson(w, map[string]string{})
}
