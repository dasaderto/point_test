package main

import (
	"github.com/gorilla/mux"
	"news/handlers"
)

func MakeRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", handlers.MainHandler).Methods("GET")
	router.HandleFunc("/resources/list", handlers.ResourcesListHandler).Methods("GET", "POST")
	router.HandleFunc("/resources/create", handlers.ResourcesCreateViewHandler).Methods("GET")
	router.HandleFunc("/resources", handlers.ResourcesCreateHandler).Methods("POST")
	router.HandleFunc("/resources/destroy", handlers.DestroyResourceHandler).Methods("POST")
	router.HandleFunc("/actualize/news", handlers.ActualizeNewsHandler).Methods("GET")

	return router
}
