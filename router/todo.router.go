package router

import (
	"github.com/gorilla/mux"
	"github.com/sumit-verma-trackier/todo/controller"
)

func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/todo", controller.AddTodoRecordController).Methods("POST")
	router.HandleFunc("/api/todo", controller.GetAllTodoController).Methods("GET")
	router.HandleFunc("/api/todo/{id}", controller.GetTodoByIdController).Methods("GET")
	router.HandleFunc("/api/todo/{id}", controller.DeleteTodoRecordController).Methods("DELETE")
	router.HandleFunc("/api/todo/{id}", controller.UpdateTodoRecordController).Methods("PUT")

	return router

}
