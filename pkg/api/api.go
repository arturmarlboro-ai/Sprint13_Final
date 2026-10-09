package api

import (
	"github.com/go-chi/chi/v5"
)

const DateFormat = "20060102"

var Router = chi.NewRouter()

func Init() {

	Router.Get("/tasks", tasksHandler)
	Router.Get("/nextdate", nextDayHandler)
	Router.Post("/task", AddTaskHandler)
	Router.Get("/task", editTaskHandler)
	Router.Put("/task", putTaskHandler)
	Router.Post("/task/done", doneTaskHandler)

}
