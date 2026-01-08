package controller

import (
	"encoding/json"
	"github.com/haqury/user-service/internal/service"
	"net/http"
)

type Controller struct {
	services *service.Services
}

func New(services *service.Services) *Controller {
	return &Controller{services: services}
}

func (c *Controller) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", c.HealthCheck)
	mux.HandleFunc("/api/v1/users", c.handleUsers)
	return mux
}

func (c *Controller) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"service": "user-service",
	})
}

func (c *Controller) handleUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message": "Users endpoint",
		"method":  r.Method,
	}

	if r.Method == "GET" {
		users, total, _ := c.services.User.ListUsers(r.Context(), 1, 10, "")
		response["users"] = users
		response["total"] = total
	}

	json.NewEncoder(w).Encode(response)
}

// Заглушки для остальных методов
func (c *Controller) ListUsers(w http.ResponseWriter, r *http.Request)             {}
func (c *Controller) GetUser(w http.ResponseWriter, r *http.Request)               {}
func (c *Controller) CreateUser(w http.ResponseWriter, r *http.Request)            {}
func (c *Controller) UpdateUser(w http.ResponseWriter, r *http.Request)            {}
func (c *Controller) DeleteUser(w http.ResponseWriter, r *http.Request)            {}
func (c *Controller) Login(w http.ResponseWriter, r *http.Request)                 {}
func (c *Controller) ValidateToken(w http.ResponseWriter, r *http.Request)         {}
func (c *Controller) Logout(w http.ResponseWriter, r *http.Request)                {}
func (c *Controller) GetStreamingConfig(w http.ResponseWriter, r *http.Request)    {}
func (c *Controller) UpdateStreamingConfig(w http.ResponseWriter, r *http.Request) {}
func (c *Controller) GetUserStats(w http.ResponseWriter, r *http.Request)          {}
