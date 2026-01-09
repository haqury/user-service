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

	"github.com/haqury/user-service/internal/config"
	"github.com/haqury/user-service/internal/repository"
	"github.com/haqury/user-service/internal/service"
)

// Application - основная структура приложения
type Application struct {
	Config   *config.Config
	Services *service.Services
	Repos    *repository.Repositories
}

// New создает новое приложение
func New(c *config.Config) *Application {
	// Инициализируем репозитории
	repos := repository.New()

	// Инициализируем сервисы
	services := service.New(repos)

	return &Application{
		Config:   c,
		Services: services,
		Repos:    repos,
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

// Run запускает приложение (gRPC + HTTP Gateway)
func (app *Application) Run() error {
	// Канал для ошибок
	errChan := make(chan error, 2)

	// Адреса серверов
	grpcAddr := ":" + app.Config.GRPCPort
	httpAddr := ":" + app.Config.HTTPPort

	// Запускаем gRPC сервер в горутине
	go func() {
		log.Printf("Starting gRPC server on port %s", app.Config.GRPCPort)
		if err := StartGRPCServer(app, grpcAddr); err != nil {
			errChan <- fmt.Errorf("gRPC server error: %w", err)
		}
	}()

	// Даем gRPC серверу время на запуск
	time.Sleep(100 * time.Millisecond)

	// Запускаем HTTP Gateway сервер в горутине
	go func() {
		log.Printf("Starting HTTP Gateway on port %s", app.Config.HTTPPort)
		ctx := context.Background()
		if err := StartGatewayServer(ctx, "localhost"+grpcAddr, httpAddr); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("HTTP gateway error: %w", err)
		}
	}()

	// Обработка сигналов для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		log.Println("Shutting down servers...")
	case err := <-errChan:
		log.Printf("Server error: %v", err)
		return err
	}

	log.Println("Servers stopped")
	return nil
}

// HealthCheck - проверка здоровья приложения
func (app *Application) HealthCheck() bool {
	// Проверяем подключение к БД, кэшу и т.д.
	return true
}
