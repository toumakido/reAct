package main

import (
	"log"

	"github.com/example/api-server/internal/handler"
	"github.com/example/api-server/pkg/middleware"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Apply middleware
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.AuthMiddleware)

	// User routes
	r.HandleFunc("/api/users", handler.GetUsers).Methods("GET")
	r.HandleFunc("/api/users/{id}", handler.GetUser).Methods("GET")
	r.HandleFunc("/api/users", handler.CreateUser).Methods("POST")

	// Product routes
	r.HandleFunc("/api/products", handler.GetProducts).Methods("GET")
	r.HandleFunc("/api/products/{id}", handler.GetProduct).Methods("GET")
	r.HandleFunc("/api/products", handler.CreateProduct).Methods("POST")
	r.HandleFunc("/api/products/{id}/stock", handler.UpdateProductStock).Methods("PATCH")

	// Create and start server with graceful shutdown
	server := NewServer(":8080", r)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
