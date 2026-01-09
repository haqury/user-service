package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/haqury/user-service/internal/models"
	"github.com/haqury/user-service/internal/repository"
)

type RoutingService interface {
	// SelectVideoService выбирает оптимальный video-service для пользователя
	SelectVideoService(ctx context.Context, user *models.User) (*models.VideoServiceInstance, error)

	// AssignInstanceToClient назначает video-service инстанс клиенту
	AssignInstanceToClient(ctx context.Context, clientID string, instanceID string) error

	// GetStreamingConfigForClient получает конфигурацию стриминга для клиента
	GetStreamingConfigForClient(ctx context.Context, userID, clientID string) (*models.StreamingConfig, error)
}

type routingService struct {
	repos *repository.Repositories
}

func NewRoutingService(repos *repository.Repositories) RoutingService {
	return &routingService{repos: repos}
}

// SelectVideoService выбирает оптимальный video-service инстанс для пользователя
// Логика выбора:
// 1. Фильтр по региону пользователя
// 2. Фильтр по тарифному плану
// 3. Сортировка по приоритету (premium инстансы имеют выше приоритет)
// 4. Выбор инстанса с наименьшей нагрузкой
func (s *routingService) SelectVideoService(ctx context.Context, user *models.User) (*models.VideoServiceInstance, error) {
	// Получаем инстансы по региону и тарифу
	instances, err := s.repos.VideoServiceInstance.GetByRegionAndTier(ctx, user.Region, user.SubscriptionTier)
	if err != nil {
		return nil, fmt.Errorf("failed to get instances: %w", err)
	}

	if len(instances) == 0 {
		// Если нет подходящих инстансов в регионе пользователя, пробуем default регион
		instances, err = s.repos.VideoServiceInstance.GetByRegionAndTier(ctx, "default", user.SubscriptionTier)
		if err != nil {
			return nil, fmt.Errorf("failed to get default instances: %w", err)
		}
	}

	if len(instances) == 0 {
		return nil, errors.New("no available video service instances")
	}

	// Возвращаем первый инстанс (уже отсортирован по priority DESC, current_load ASC)
	return instances[0], nil
}

// AssignInstanceToClient назначает video-service инстанс клиенту
func (s *routingService) AssignInstanceToClient(ctx context.Context, clientID string, instanceID string) error {
	return s.repos.UserClient.AssignInstance(ctx, clientID, instanceID)
}

// GetStreamingConfigForClient получает конфигурацию стриминга для клиента
func (s *routingService) GetStreamingConfigForClient(ctx context.Context, userID, clientID string) (*models.StreamingConfig, error) {
	// Проверяем, есть ли уже назначенный инстанс для этого клиента
	userClient, err := s.repos.UserClient.GetByClientID(ctx, clientID)
	var instance *models.VideoServiceInstance

	if err != nil || userClient == nil || userClient.AssignedInstanceID == nil {
		// Если клиент еще не зарегистрирован или нет назначенного инстанса,
		// нужно получить пользователя и выбрать инстанс
		user, err := s.getUserByID(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to get user: %w", err)
		}

		// Выбираем оптимальный инстанс
		instance, err = s.SelectVideoService(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("failed to select video service: %w", err)
		}

		// Если клиент существует, обновляем его назначение
		if userClient != nil {
			err = s.repos.UserClient.AssignInstance(ctx, clientID, instance.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to assign instance: %w", err)
			}
		} else {
			// Создаем новую запись клиента
			userClient = &models.UserClient{
				UserID:             userID,
				ClientID:           clientID,
				AssignedInstanceID: &instance.ID,
				IsActive:           true,
			}
			err = s.repos.UserClient.Create(ctx, userClient)
			if err != nil {
				return nil, fmt.Errorf("failed to create user client: %w", err)
			}
		}
	} else {
		// Получаем инстанс по ID
		instance, err = s.repos.VideoServiceInstance.GetByID(ctx, *userClient.AssignedInstanceID)
		if err != nil {
			return nil, fmt.Errorf("failed to get instance: %w", err)
		}
	}

	// Формируем конфигурацию стриминга
	config := &models.StreamingConfig{
		ServerURL:      instance.ServerURL,
		ServerPort:     instance.ServerPort,
		UseSSL:         instance.UseSSL,
		APIKey:         s.generateAPIKey(userID, clientID), // TODO: implement API key generation
		StreamEndpoint: instance.StreamEndpoint,
		MaxBitrate:     instance.MaxBitrate,
		MaxResolution:  instance.MaxResolution,
		Codec:          instance.Codec,
	}

	return config, nil
}

// getUserByID - вспомогательная функция для получения пользователя
// TODO: implement properly через UserRepository
func (s *routingService) getUserByID(ctx context.Context, userID string) (*models.User, error) {
	// Временная заглушка - нужно будет реализовать через UserRepository
	return &models.User{
		ID:               userID,
		IsActive:         true,
		SubscriptionTier: "free",
		Region:           "default",
	}, nil
}

// generateAPIKey генерирует API ключ для клиента
// TODO: implement proper API key generation with JWT or similar
func (s *routingService) generateAPIKey(userID, clientID string) string {
	return fmt.Sprintf("api_key_%s_%s", userID, clientID)
}
