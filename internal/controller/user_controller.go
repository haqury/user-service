package controller

import (
	"context"
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type StreamingConfig struct {
	ServerURL      string   `json:"server_url"`
	ServerPort     int      `json:"server_port"`
	APIKey         string   `json:"api_key"`
	StreamEndpoint string   `json:"stream_endpoint"`
	MaxBitrate     int      `json:"max_bitrate"`
	MaxResolution  int      `json:"max_resolution"`
	AllowedCodecs  []string `json:"allowed_codecs"`
}

type UserController struct {
	// В реальном приложении здесь был бы репозиторий
}

func NewUserController() *UserController {
	return &UserController{}
}

func (c *UserController) GetUser(ctx context.Context, userID string) (*User, error) {
	// Тестовые данные
	return &User{
		ID:        userID,
		Username:  "user_" + userID,
		Email:     "user_" + userID + "@example.com",
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (c *UserController) GetStreamingConfig(ctx context.Context, userID string) (*StreamingConfig, error) {
	return &StreamingConfig{
		ServerURL:      "http://localhost:8082",
		ServerPort:     8082,
		APIKey:         "video_api_key_" + userID,
		StreamEndpoint: "/api/v1/video/stream",
		MaxBitrate:     5000,
		MaxResolution:  1080,
		AllowedCodecs:  []string{"h264", "h265"},
	}, nil
}

func (c *UserController) Login(ctx context.Context, username, password string) (string, *User, error) {
	// Тестовая аутентификация
	token := "jwt_token_" + time.Now().Format("20060102150405")

	user := &User{
		ID:        "auth_user_001",
		Username:  username,
		Email:     username + "@example.com",
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return token, user, nil
}
