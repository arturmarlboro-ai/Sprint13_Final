package api

import (
	"github.com/go-chi/chi/v5"
)

var Router = chi.NewRouter()

func Init() {
	Router.Post("/signin", logHandler)
	Router.Get("/nextdate", nextDayHandler)

	Router.Group(func(r chi.Router) {
		r.Use(auth)
		r.Get("/task", editTaskHandler)
		r.Get("/tasks", tasksHandler)
		r.Post("/task", AddTaskHandler)
		r.Put("/task", putTaskHandler)
		r.Post("/task/done", doneTaskHandler)
		r.Delete("/task", deleteTaskHandler)
	})

}
