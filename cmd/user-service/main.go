package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"user-service/internal/app"
	"user-service/internal/config"

	"go.uber.org/zap"
)

func main() {
	// Определяем флаги
	configPath := flag.String("config", "config.yaml", "Path to config file")
	grpcPort := flag.String("grpc-port", "9091", "gRPC server port")
	flag.Parse()

	// Инициализируем логгер
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting user service",
		zap.String("config", *configPath),
		zap.String("grpc_port", *grpcPort),
	)

	// Загрузка конфигурации
	cfg, err := config.Load(*configPath)
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	// Создание DSN строки для подключения к БД
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)
	// Создание приложения
	application := app.New(&app.Config{
		HTTPPort: fmt.Sprintf("%d", cfg.Server.Port),
		GRPCPort: *grpcPort,
		Env:      cfg.Server.Mode,
		Database: struct {
			DSN string
		}{
			DSN: dsn,
		},
	})
	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-quit
		logger.Info("Received signal, shutting down", zap.String("signal", sig.String()))
		// Здесь нужно добавить shutdown для application
		os.Exit(0)
	}()

	// Запуск приложения
	if err := application.Run(); err != nil {
		logger.Fatal("Application error", zap.Error(err))
	}

	logger.Info("Application stopped")
}
