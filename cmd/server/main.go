package main

import (
	"golang-abac-demo/internal/config"
	"golang-abac-demo/internal/controllers"
	"golang-abac-demo/internal/middlewares"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize the router
	r := mux.NewRouter()

	// Log all requests
	r.Use(middlewares.LoggingMiddleware)

	// Public Routes
	r.HandleFunc("/login", controllers.Login).Methods("POST")

	// Private Routes
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middlewares.AuthMiddleware)
	api.HandleFunc("/documents", controllers.UploadDocument).Methods("POST")
	api.HandleFunc("/documents/{id}", controllers.ViewDocument).Methods("GET")
	api.HandleFunc("/documents/{id}", controllers.EditDocument).Methods("PUT")
	api.HandleFunc("/documents/{id}", controllers.DeleteDocument).Methods("DELETE")

	// Start the server
	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
