package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"userservice/internal/service"
	"userservice/pkg/gen"
)

type UserHandler struct {
	logger      *zap.Logger
	userService service.UserService
}

func NewUserHandler(logger *zap.Logger, userService service.UserService) *UserHandler {
	return &UserHandler{
		logger:      logger,
		userService: userService,
	}
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/health":
		h.healthCheck(w, r)
	case "/api/v1/user":
		h.getUser(w, r)
	case "/api/v1/login":
		h.login(w, r)
	case "/api/v1/streaming/config":
		h.getStreamingConfig(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *UserHandler) healthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]interface{}{
		"status":  "ok",
		"service": "user-service",
		"version": "1.0.0",
		"time":    time.Now().Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) getUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUser(r.Context(), &pb.GetUserRequest{UserId: userID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	loginReq := &pb.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}

	response, err := h.userService.Login(r.Context(), loginReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) getStreamingConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	clientID := r.URL.Query().Get("client_id")

	if userID == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	req := &pb.GetStreamingConfigRequest{
		UserId:   userID,
		ClientId: clientID,
	}

	config, err := h.userService.GetStreamingConfig(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}
