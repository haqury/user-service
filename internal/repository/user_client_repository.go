package repository

import (
	"context"
	"database/sql"

	"github.com/haqury/user-service/internal/models"
)

type UserClientRepository interface {
	GetByClientID(ctx context.Context, clientID string) (*models.UserClient, error)
	GetByUserID(ctx context.Context, userID string) ([]*models.UserClient, error)
	Create(ctx context.Context, userClient *models.UserClient) error
	Update(ctx context.Context, userClient *models.UserClient) error
	AssignInstance(ctx context.Context, clientID, instanceID string) error
	UpdateLastSeen(ctx context.Context, clientID string) error
	Delete(ctx context.Context, id string) error
}

type userClientRepository struct {
	db *sql.DB
}

func NewUserClientRepository(db *sql.DB) UserClientRepository {
	return &userClientRepository{db: db}
}

func (r *userClientRepository) GetByClientID(ctx context.Context, clientID string) (*models.UserClient, error) {
	query := `
		SELECT id, user_id, client_id, client_info, assigned_instance_id,
		       is_active, created_at, updated_at, last_seen
		FROM user_clients
		WHERE client_id = $1 AND is_active = true
	`

	var uc models.UserClient
	err := r.db.QueryRowContext(ctx, query, clientID).Scan(
		&uc.ID, &uc.UserID, &uc.ClientID, &uc.ClientInfo, &uc.AssignedInstanceID,
		&uc.IsActive, &uc.CreatedAt, &uc.UpdatedAt, &uc.LastSeen,
	)
	if err != nil {
		return nil, err
	}

	return &uc, nil
}

func (r *userClientRepository) GetByUserID(ctx context.Context, userID string) ([]*models.UserClient, error) {
	// TODO: implement
	return nil, nil
}

func (r *userClientRepository) Create(ctx context.Context, userClient *models.UserClient) error {
	query := `
		INSERT INTO user_clients (user_id, client_id, client_info, assigned_instance_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at, last_seen
	`

	return r.db.QueryRowContext(
		ctx, query,
		userClient.UserID, userClient.ClientID, userClient.ClientInfo, userClient.AssignedInstanceID,
	).Scan(&userClient.ID, &userClient.CreatedAt, &userClient.UpdatedAt, &userClient.LastSeen)
}

func (r *userClientRepository) Update(ctx context.Context, userClient *models.UserClient) error {
	// TODO: implement
	return nil
}

func (r *userClientRepository) AssignInstance(ctx context.Context, clientID, instanceID string) error {
	query := `UPDATE user_clients SET assigned_instance_id = $1 WHERE client_id = $2`
	_, err := r.db.ExecContext(ctx, query, instanceID, clientID)
	return err
}

func (r *userClientRepository) UpdateLastSeen(ctx context.Context, clientID string) error {
	query := `UPDATE user_clients SET last_seen = CURRENT_TIMESTAMP WHERE client_id = $1`
	_, err := r.db.ExecContext(ctx, query, clientID)
	return err
}

func (r *userClientRepository) Delete(ctx context.Context, id string) error {
	// TODO: implement
	return nil
}
