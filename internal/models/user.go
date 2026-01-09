package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/lib/pq"
)

type User struct {
	ID               string         `db:"id" json:"id"`
	Username         string         `db:"username" json:"username"`
	Email            string         `db:"email" json:"email"`
	Phone            string         `db:"phone" json:"phone"`
	PasswordHash     string         `db:"password_hash" json:"-"`
	Status           string         `db:"status" json:"status"`
	IsActive         bool           `db:"is_active" json:"is_active"`
	Roles            pq.StringArray `db:"roles" json:"roles"`
	SubscriptionTier string         `db:"subscription_tier" json:"subscription_tier"`
	Region           string         `db:"region" json:"region"`
	Settings         JSONB          `db:"settings" json:"settings"`
	StreamingConfig  JSONB          `db:"streaming_config" json:"streaming_config"`
	Stats            JSONB          `db:"stats" json:"stats"`
	Metadata         JSONB          `db:"metadata" json:"metadata"`
	CreatedAt        time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time      `db:"updated_at" json:"updated_at"`
	LastLogin        *time.Time     `db:"last_login" json:"last_login,omitempty"`
	LastActivity     *time.Time     `db:"last_activity" json:"last_activity,omitempty"`
}

// JSONB is a custom type for JSONB fields
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	result := make(JSONB)
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}

	*j = result
	return nil
}

type UserSettings struct {
	DefaultQuality       string `json:"default_quality"`
	MaxParallelStreams   int32  `json:"max_parallel_streams"`
	AutoStartRecording   bool   `json:"auto_start_recording"`
	NotificationsEnabled bool   `json:"notifications_enabled"`
	Timezone             string `json:"timezone"`
	Language             string `json:"language"`
}

type StreamingConfig struct {
	ServerURL      string `json:"server_url"`
	ServerPort     int32  `json:"server_port"`
	UseSSL         bool   `json:"use_ssl"`
	APIKey         string `json:"api_key"`
	StreamEndpoint string `json:"stream_endpoint"`
	MaxBitrate     int32  `json:"max_bitrate"`
	MaxResolution  int32  `json:"max_resolution"`
	Codec          string `json:"codec"`
}

type UserStats struct {
	TotalStreams     int64 `json:"total_streams"`
	TotalDuration    int64 `json:"total_duration"`
	TotalStorageUsed int64 `json:"total_storage_used"`
	CurrentStreams   int64 `json:"current_streams"`
	SuccessfulLogins int32 `json:"successful_logins"`
	FailedLogins     int32 `json:"failed_logins"`
	LastLogin        int64 `json:"last_login"`
	LastActivity     int64 `json:"last_activity"`
}
