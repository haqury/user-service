package main

import (
	"flag"
	"github.com/haqury/user-service/internal/app"
	"log"

	"go.uber.org/zap"
)

func main() {
	// Определяем флаги
	configPath := flag.String("config", "", "Path to config file (optional, uses env vars if not provided)")
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

	// Создаем приложение с конфигурацией
	application, err := app.NewWithConfig(*configPath, *grpcPort)
	if err != nil {
		logger.Fatal("Failed to create application", zap.Error(err))
	}

	logger.Info("Application created successfully")
	// Запуск приложения
	if err := application.Run(); err != nil {
		logger.Fatal("Application error", zap.Error(err))
	}

	logger.Info("Application stopped")
}
