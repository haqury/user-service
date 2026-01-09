package models

import (
	"time"

	"github.com/lib/pq"
)

type VideoServiceInstance struct {
	ID              string         `db:"id" json:"id"`
	Name            string         `db:"name" json:"name"`
	ServerURL       string         `db:"server_url" json:"server_url"`
	ServerPort      int32          `db:"server_port" json:"server_port"`
	UseSSL          bool           `db:"use_ssl" json:"use_ssl"`
	StreamEndpoint  string         `db:"stream_endpoint" json:"stream_endpoint"`
	Region          string         `db:"region" json:"region"`
	Priority        int32          `db:"priority" json:"priority"`
	MaxCapacity     int32          `db:"max_capacity" json:"max_capacity"`
	CurrentLoad     int32          `db:"current_load" json:"current_load"`
	HealthStatus    string         `db:"health_status" json:"health_status"`
	AllowedTiers    pq.StringArray `db:"allowed_tiers" json:"allowed_tiers"`
	MaxBitrate      int32          `db:"max_bitrate" json:"max_bitrate"`
	MaxResolution   int32          `db:"max_resolution" json:"max_resolution"`
	Codec           string         `db:"codec" json:"codec"`
	Metadata        JSONB          `db:"metadata" json:"metadata"`
	IsActive        bool           `db:"is_active" json:"is_active"`
	CreatedAt       time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time      `db:"updated_at" json:"updated_at"`
	LastHealthCheck *time.Time     `db:"last_health_check" json:"last_health_check,omitempty"`
}
