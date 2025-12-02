package main

import (
	"log"
	"net/http"

	"github.com/example/api-server/internal/handler"
	"github.com/example/api-server/pkg/middleware"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Apply middleware
	r.Use(middleware.AuthMiddleware)

	// User routes
	r.HandleFunc("/api/users", handler.GetUsers).Methods("GET")
	r.HandleFunc("/api/users/{id}", handler.GetUser).Methods("GET")
	r.HandleFunc("/api/users", handler.CreateUser).Methods("POST")

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
