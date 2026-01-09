package repository

import (
	"context"
	"database/sql"

	"github.com/haqury/user-service/internal/models"
)

type VideoServiceInstanceRepository interface {
	GetByID(ctx context.Context, id string) (*models.VideoServiceInstance, error)
	GetByName(ctx context.Context, name string) (*models.VideoServiceInstance, error)
	GetActiveInstances(ctx context.Context) ([]*models.VideoServiceInstance, error)
	GetByRegionAndTier(ctx context.Context, region, tier string) ([]*models.VideoServiceInstance, error)
	UpdateLoad(ctx context.Context, id string, load int32) error
	Create(ctx context.Context, instance *models.VideoServiceInstance) error
	Update(ctx context.Context, instance *models.VideoServiceInstance) error
	Delete(ctx context.Context, id string) error
}

type videoServiceInstanceRepository struct {
	db *sql.DB
}

func NewVideoServiceInstanceRepository(db *sql.DB) VideoServiceInstanceRepository {
	return &videoServiceInstanceRepository{db: db}
}

func (r *videoServiceInstanceRepository) GetByID(ctx context.Context, id string) (*models.VideoServiceInstance, error) {
	// TODO: implement
	return nil, nil
}

func (r *videoServiceInstanceRepository) GetByName(ctx context.Context, name string) (*models.VideoServiceInstance, error) {
	// TODO: implement
	return nil, nil
}

func (r *videoServiceInstanceRepository) GetActiveInstances(ctx context.Context) ([]*models.VideoServiceInstance, error) {
	// TODO: implement
	return nil, nil
}

func (r *videoServiceInstanceRepository) GetByRegionAndTier(ctx context.Context, region, tier string) ([]*models.VideoServiceInstance, error) {
	query := `
		SELECT id, name, server_url, server_port, use_ssl, stream_endpoint, region,
		       priority, max_capacity, current_load, health_status, allowed_tiers,
		       max_bitrate, max_resolution, codec, metadata, is_active,
		       created_at, updated_at, last_health_check
		FROM video_service_instances
		WHERE is_active = true
		  AND health_status = 'healthy'
		  AND region = $1
		  AND $2 = ANY(allowed_tiers)
		  AND current_load < max_capacity
		ORDER BY priority DESC, current_load ASC
	`

	// TODO: implement query execution
	_ = query
	return nil, nil
}

func (r *videoServiceInstanceRepository) UpdateLoad(ctx context.Context, id string, load int32) error {
	query := `UPDATE video_service_instances SET current_load = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, load, id)
	return err
}

func (r *videoServiceInstanceRepository) Create(ctx context.Context, instance *models.VideoServiceInstance) error {
	// TODO: implement
	return nil
}

func (r *videoServiceInstanceRepository) Update(ctx context.Context, instance *models.VideoServiceInstance) error {
	// TODO: implement
	return nil
}

func (r *videoServiceInstanceRepository) Delete(ctx context.Context, id string) error {
	// TODO: implement
	return nil
}
