package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/haqury/user-service/pkg/gen"
)

// StartGatewayServer запускает HTTP Gateway сервер, который проксирует запросы в gRPC
func StartGatewayServer(ctx context.Context, grpcAddr, httpAddr string) error {
	// Создаем mux для gRPC-Gateway
	mux := runtime.NewServeMux()

	// Опции для подключения к gRPC серверу
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// Регистрируем UserService handler
	err := pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts)
	if err != nil {
		return fmt.Errorf("failed to register gateway: %w", err)
	}

	// Добавляем health check endpoint
	healthMux := http.NewServeMux()
	healthMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"user-service"}`))
	})

	// Все остальные запросы идут в gRPC-Gateway
	healthMux.Handle("/", mux)

	// Запускаем HTTP сервер
	server := &http.Server{
		Addr:    httpAddr,
		Handler: allowCORS(healthMux),
	}

	return server.ListenAndServe()
}

// allowCORS добавляет CORS headers
func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept", "Authorization"}
	w.Header().Set("Access-Control-Allow-Headers", joinHeaders(headers))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", joinHeaders(methods))
}

func joinHeaders(headers []string) string {
	result := ""
	for i, h := range headers {
		if i > 0 {
			result += ", "
		}
		result += h
	}
	return result
}
