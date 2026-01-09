package handler

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/haqury/helpy"
	"github.com/haqury/user-service/internal/service"
	pb "github.com/haqury/user-service/pkg/gen"
)

// UserServiceServer реализует gRPC сервер для UserService
type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
	userService    service.UserService
	routingService service.RoutingService
}

func NewUserServiceServer(userService service.UserService, routingService service.RoutingService) *UserServiceServer {
	return &UserServiceServer{
		userService:    userService,
		routingService: routingService,
	}
}

// GetUser получает пользователя по ID
func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	user, err := s.userService.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	// TODO: преобразовать user из interface{} в *pb.User
	_ = user
	return &pb.User{
		Id:       req.UserId,
		Username: "test",
		Email:    "test@example.com",
	}, nil
}

// GetUserByUsername получает пользователя по имени
func (s *UserServiceServer) GetUserByUsername(ctx context.Context, req *pb.GetUserByUsernameRequest) (*pb.User, error) {
	user, err := s.userService.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	// TODO: преобразовать user из interface{} в *pb.User
	_ = user
	return &pb.User{
		Username: req.Username,
	}, nil
}

// GetUserByClientId получает информацию о пользователе по client_id
func (s *UserServiceServer) GetUserByClientId(ctx context.Context, req *pb.GetUserByClientIdRequest) (*pb.GetUserByClientIdResponse, error) {
	userInfo, err := s.userService.GetUserByClientID(ctx, req.ClientId)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserByClientIdResponse{
		UserId:   userInfo.UserID,
		Username: userInfo.Username,
		IsActive: userInfo.IsActive,
		Roles:    userInfo.Roles,
	}, nil
}

// CreateUser создает нового пользователя
func (s *UserServiceServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	// TODO: implement
	return nil, nil
}

// UpdateUser обновляет пользователя
func (s *UserServiceServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	// TODO: implement
	return nil, nil
}

// DeleteUser удаляет пользователя
func (s *UserServiceServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*helpy.ApiResponse, error) {
	// TODO: implement
	return nil, nil
}

// ListUsers возвращает список пользователей
func (s *UserServiceServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	// TODO: implement
	return nil, nil
}

// Login аутентифицирует пользователя
func (s *UserServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// TODO: implement
	return nil, nil
}

// ValidateToken проверяет токен
func (s *UserServiceServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	// TODO: implement
	return nil, nil
}

// Logout выполняет выход пользователя
func (s *UserServiceServer) Logout(ctx context.Context, req *emptypb.Empty) (*helpy.ApiResponse, error) {
	// TODO: implement
	return nil, nil
}

// GetStreamingConfig получает конфигурацию стриминга для пользователя
func (s *UserServiceServer) GetStreamingConfig(ctx context.Context, req *pb.GetStreamingConfigRequest) (*pb.User_StreamingConfig, error) {
	config, err := s.routingService.GetStreamingConfigForClient(ctx, req.UserId, req.ClientId)
	if err != nil {
		return nil, err
	}

	return &pb.User_StreamingConfig{
		ServerUrl:      config.ServerURL,
		ServerPort:     config.ServerPort,
		UseSsl:         config.UseSSL,
		ApiKey:         config.APIKey,
		StreamEndpoint: config.StreamEndpoint,
		MaxBitrate:     config.MaxBitrate,
		MaxResolution:  config.MaxResolution,
		Codec:          config.Codec,
	}, nil
}

// UpdateStreamingConfig обновляет конфигурацию стриминга
func (s *UserServiceServer) UpdateStreamingConfig(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User_StreamingConfig, error) {
	// TODO: implement
	return nil, nil
}

// UpdateUserStats обновляет статистику пользователя
func (s *UserServiceServer) UpdateUserStats(ctx context.Context, req *pb.UpdateUserRequest) (*helpy.ApiResponse, error) {
	// TODO: implement
	return nil, nil
}

// GetUserStats получает статистику пользователя
func (s *UserServiceServer) GetUserStats(ctx context.Context, req *pb.GetUserRequest) (*pb.User_UserStats, error) {
	// TODO: implement
	return nil, nil
}
