package app

import (
	"user-service/internal/controller"

	"github.com/gorilla/mux"
)

func RegisterRoutes(ctrl *controller.Controller) *mux.Router {
	router := mux.NewRouter()

	// Health check
	router.HandleFunc("/health", ctrl.HealthCheck).Methods("GET")

	// API v1
	api := router.PathPrefix("/api/v1").Subrouter()

	// User endpoints
	api.HandleFunc("/users", ctrl.ListUsers).Methods("GET")
	api.HandleFunc("/users", ctrl.CreateUser).Methods("POST")
	api.HandleFunc("/users/{id}", ctrl.GetUser).Methods("GET")
	api.HandleFunc("/users/{id}", ctrl.UpdateUser).Methods("PUT")
	api.HandleFunc("/users/{id}", ctrl.DeleteUser).Methods("DELETE")

	// Auth endpoints
	api.HandleFunc("/auth/login", ctrl.Login).Methods("POST")
	api.HandleFunc("/auth/validate", ctrl.ValidateToken).Methods("POST")
	api.HandleFunc("/auth/logout", ctrl.Logout).Methods("POST")

	// Streaming config
	api.HandleFunc("/users/{id}/streaming-config", ctrl.GetStreamingConfig).Methods("GET")
	api.HandleFunc("/users/{id}/streaming-config", ctrl.UpdateStreamingConfig).Methods("PUT")

	// User stats
	api.HandleFunc("/users/{id}/stats", ctrl.GetUserStats).Methods("GET")

	return router
}
