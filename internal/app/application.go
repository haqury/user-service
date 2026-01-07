package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user-service/internal/config"

	"user-service/internal/controller"
	"user-service/internal/repository"
	"user-service/internal/service"
)

// Application - основная структура приложения
type Application struct {
	Config     *config.Config
	HTTPServer *http.Server
	GRPCServer interface{} // Можно добавить позже
	Services   *service.Services
	Repos      *repository.Repositories
	Controller *controller.Controller
}

// New создает новое приложение
func New(c *config.Config) *Application {
	// Инициализируем репозитории
	repos := repository.New()

	// Инициализируем сервисы
	services := service.New(repos)

	// Инициализируем контроллер
	ctrl := controller.New(services)

	// Настраиваем HTTP сервер
	mux := ctrl.RegisterRoutes()

	server := &http.Server{
		Addr:         ":" + c.HTTPPort,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Application{
		Config:     c,
		HTTPServer: server,
		Services:   services,
		Repos:      repos,
		Controller: ctrl,
	}
}

// NewWithConfig создает приложение с конфигурацией из файла или env
func NewWithConfig(cfgPath string, grpcPort string) (*Application, error) {
	c, err := config.NewConfig(cfgPath, grpcPort)
	if err != nil {
		return nil, fmt.Errorf("failed to create config: %w", err)
	}

	return New(c), nil
}

// Run запускает приложение
func (app *Application) Run() error {
	// Канал для ошибок
	errChan := make(chan error, 1)

	// Запускаем HTTP сервер в горутине
	go func() {
		log.Printf("Starting HTTP server on port %s", app.Config.HTTPPort)
		if err := app.HTTPServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Обработка сигналов для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		log.Println("Shutting down server...")
	case err := <-errChan:
		log.Printf("Server error: %v", err)
	}

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.HTTPServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
		return err
	}

	log.Println("Server stopped")
	return nil
}

// HealthCheck - проверка здоровья приложения
func (app *Application) HealthCheck() bool {
	// Проверяем подключение к БД, кэшу и т.д.
	return true
}
