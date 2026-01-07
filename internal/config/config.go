package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type PreConfig struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	Log      LogConfig      `yaml:"log"`
}

type ServerConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	Mode         string        `yaml:"mode"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	SSLMode  string `yaml:"ssl_mode"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type JWTConfig struct {
	Secret     string        `yaml:"secret"`
	Expiration time.Duration `yaml:"expiration"`
}

type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type ServicesConfig struct {
	APIGatewayURL   string `yaml:"api_gateway_url"`
	VideoServiceURL string `yaml:"video_service_url"`
}

// Config - конфигурация приложения
type Config struct {
	HTTPPort string
	GRPCPort string
	Env      string
	Database struct {
		DSN string
	}
}

// NewConfig создает конфигурацию приложения
func NewConfig(cfgPath string, grpcPort string) (*Config, error) {
	var appConfig *PreConfig
	var err error

	if cfgPath != "" {
		appConfig, err = Load(cfgPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load config: %w", err)
		}
	} else {
		appConfig = LoadFromEnv()
	}

	// Создание DSN строки для подключения к БД
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		appConfig.Database.User,
		appConfig.Database.Password,
		appConfig.Database.Host,
		appConfig.Database.Port,
		appConfig.Database.Name,
		appConfig.Database.SSLMode,
	)

	return &Config{
		HTTPPort: fmt.Sprintf("%d", appConfig.Server.Port),
		GRPCPort: grpcPort,
		Env:      appConfig.Server.Mode,
		Database: struct {
			DSN string
		}{
			DSN: dsn,
		},
	}, nil
}

// overrideFromEnv переопределяет значения из переменных окружения
func overrideFromEnv(cfg *PreConfig) *PreConfig {
	// Database
	if host := os.Getenv("DB_HOST"); host != "" {
		cfg.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Database.Port = p
		}
	}
	if user := os.Getenv("DB_USER"); user != "" {
		cfg.Database.User = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		cfg.Database.Password = password
	}
	if name := os.Getenv("DB_NAME"); name != "" {
		cfg.Database.Name = name
	}
	if sslMode := os.Getenv("DB_SSLMODE"); sslMode != "" {
		cfg.Database.SSLMode = sslMode
	}

	// Redis
	if host := os.Getenv("REDIS_HOST"); host != "" {
		cfg.Redis.Host = host
	}
	if port := os.Getenv("REDIS_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Redis.Port = p
		}
	}
	if password := os.Getenv("REDIS_PASSWORD"); password != "" {
		cfg.Redis.Password = password
	}
	if db := os.Getenv("REDIS_DB"); db != "" {
		if d, err := strconv.Atoi(db); err == nil {
			cfg.Redis.DB = d
		}
	}

	// Server
	if host := os.Getenv("SERVER_HOST"); host != "" {
		cfg.Server.Host = host
	}
	if port := os.Getenv("SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Server.Port = p
		}
	}
	if mode := os.Getenv("SERVER_MODE"); mode != "" {
		cfg.Server.Mode = mode
	}

	// Log
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		cfg.Log.Level = level
	}
	if format := os.Getenv("LOG_FORMAT"); format != "" {
		cfg.Log.Format = format
	}

	return cfg
}

// Load загружает конфигурацию из файла
func Load(path string) (*PreConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		// Не возвращаем Default, возвращаем ошибку
		return nil, fmt.Errorf("failed to load config file %s: %v", path, err)
	}

	cfg := Default()
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Переопределяем значения из переменных окружения
	cfg = overrideFromEnv(cfg)

	return cfg, nil
}

// LoadFromEnv загружает конфигурацию ТОЛЬКО из переменных окружения
// Сначала пытается загрузить .env файл, затем устанавливает значения по умолчанию
// и переопределяет их значениями из переменных окружения
func LoadFromEnv() *PreConfig {
	// Пытаемся загрузить .env файл (игнорируем ошибку, если файл не найден)
	_ = godotenv.Load()

	// Начинаем с значений по умолчанию
	cfg := Default()

	// Переопределяем значения из переменных окружения
	cfg = overrideFromEnv(cfg)

	return cfg
}

// Default возвращает конфигурацию по умолчанию
func Default() *PreConfig {
	return &PreConfig{
		Server: ServerConfig{
			Host:         "0.0.0.0",
			Port:         8081,
			Mode:         "debug",
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "postgres",
			Name:     "user_service",
			SSLMode:  "disable",
		},
		Redis: RedisConfig{
			Host: "localhost",
			Port: 6379,
			DB:   0,
		},
		Log: LogConfig{
			Level:  "info",
			Format: "json",
		},
	}
}
