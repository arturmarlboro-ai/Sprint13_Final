package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"Sprint13_Final/pkg/constants"
	"Sprint13_Final/pkg/db"
)

func nextDayHandler(w http.ResponseWriter, r *http.Request) {

	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeatStr := r.URL.Query().Get("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(constants.DateFormat, nowStr)
		if err != nil {
			log.Printf("error time parse at handlers.go")
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": "time parse at handlers.go"})
			return
		}
	}

	nextDateStr, err := NextDate(now, dateStr, repeatStr)
	if err != nil {
		log.Print("error at NextDateHandler, nextDateStr")
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": "at NextDateHandler"})
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "%s", nextDateStr)

}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	LimitTasks := 50

	var tasks []*db.Task
	var err error
	search := r.URL.Query().Get("search")

	log.Print(search)

	if search == "" {
		tasks, err = db.GetTasksByLimit(r.Context(), LimitTasks)
	} else {
		tasks, err = db.GetTasksBySearch(r.Context(), search)
	}

	if err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": "failed to get tasks"})
		return
	}

	writeJson(w, http.StatusOK, TasksResp{
		Tasks: tasks,
	})
}

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Printf("PostTask json decode error: %s", err)
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}

	if task.Title == "" {
		log.Println("task.Title emtpy")
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "title is required"})
		return
	}

	if task.Date == "" {
		task.Date = time.Now().Format(constants.DateFormat)
	}

	if err := checkDate(&task); err != nil {
		log.Printf("error checkDate: %v", err)
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "invalid date"})
		return
	}

	id, err := db.AddTask(r.Context(), task)
	if err != nil {
		log.Printf("AddTask error: %v", err)
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": "failed to add task"})
		return
	}

	writeJson(w, http.StatusCreated, map[string]string{"id": id})

}

func editTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "missing id"})
		return
	}

	task, err := db.GetTaskByID(r.Context(), id)
	if err != nil {
		log.Printf("editTaskHandler error: %v", err)
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": "error getting task for id"})
		return
	}

	if task == nil {
		writeJson(w, http.StatusNotFound, map[string]string{"error": "no task for id"})
		return
	}
	log.Print(task)
	writeJson(w, http.StatusOK, task)
}

func putTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Printf("putTask json decode error: %v", err)
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	if task.ID == "" {
		log.Print("putTask missing id")
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "missing id"})
		return
	}

	if task.Title == "" {
		log.Print("putTask missing title")
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "missing title"})
		return
	}

	if err := checkDate(&task); err != nil {
		log.Printf("error checkDate: %v", err)
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "invalid date"})
		return
	}

	err := db.UpdateTask(r.Context(), task, task.ID)
	if err != nil {
		log.Printf("putTask error: %v", err)
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": "failed to update task"})
		return
	}

	writeJson(w, http.StatusOK, map[string]string{})
}

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	task, err := db.GetTaskByID(r.Context(), id)
	if err != nil {
		log.Printf("doneTaskHandler error: %v", err)
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": "failed get task by id"})
		return
	}
	if task.Repeat != "" {
		if task.Date == "" {
			task.Date = time.Now().Format(constants.DateFormat)
		}

		nextDateStr, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			log.Printf("doneTaskHandler error: %v", err)
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": "failed to calculate next date"})
			return
		}
		task.Date = nextDateStr

		err = db.UpdateTask(r.Context(), *task, task.ID)
		if err != nil {
			log.Printf("doneTaskHandler error: %v", err)
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": "failed to update task with next date"})
			return
		}
	} else {
		err = db.DeleteTask(r.Context(), id)
		if err != nil {
			log.Printf("doneTaskHandler error: %v", err)
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": "failed to delete task"})
			return
		}
	}

	writeJson(w, http.StatusOK, map[string]string{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		log.Print("missing id in deleteTaskHandler")
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "not recieved id"})
		return
	}

	err := db.DeleteTask(r.Context(), id)
	if err != nil {
		log.Printf("error DeleteTaskHandler: %v", err)
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": "failed to delete task"})
		return
	}

	writeJson(w, http.StatusOK, map[string]string{})
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	pass := os.Getenv("TODO_PASSWORD")
	if pass == "" || req.Password != pass {
		writeJson(w, http.StatusUnauthorized, map[string]string{"error": "Неверный пароль"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		//"exp": time.Now().Add(time.Hour * 8).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(pass))
	if err != nil {
		log.Printf("logHandler error signing token: %v", err)
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": "failed to generate token"})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenStr,
		Expires: time.Now().Add(8 * time.Hour),
		Path:    "/",
	})

	writeJson(w, http.StatusOK, map[string]string{"token": tokenStr})
}
