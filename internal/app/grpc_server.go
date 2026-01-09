package app

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/haqury/user-service/internal/handler"
	pb "github.com/haqury/user-service/pkg/gen"
)

// StartGRPCServer запускает gRPC сервер
func StartGRPCServer(app *Application, addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	// Создаем gRPC сервер
	grpcServer := grpc.NewServer()

	// Создаем handler
	userServiceServer := handler.NewUserServiceServer(
		app.Services.User,
		app.Services.Routing,
	)

	// Регистрируем сервис
	pb.RegisterUserServiceServer(grpcServer, userServiceServer)

	// Включаем reflection для grpcurl и подобных инструментов
	reflection.Register(grpcServer)

	log.Printf("Starting gRPC server on %s", addr)
	return grpcServer.Serve(lis)
}
