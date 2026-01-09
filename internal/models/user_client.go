package models

import (
	"time"
)

type UserClient struct {
	ID                 string    `db:"id" json:"id"`
	UserID             string    `db:"user_id" json:"user_id"`
	ClientID           string    `db:"client_id" json:"client_id"`
	ClientInfo         JSONB     `db:"client_info" json:"client_info"`
	AssignedInstanceID *string   `db:"assigned_instance_id" json:"assigned_instance_id,omitempty"`
	IsActive           bool      `db:"is_active" json:"is_active"`
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time `db:"updated_at" json:"updated_at"`
	LastSeen           time.Time `db:"last_seen" json:"last_seen"`
}
