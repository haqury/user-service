package main

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "github.com/uptrace/bun/driver/pgdriver" // PostgreSQL драйвер
)

// Config представляет минимальную конфигурацию для миграций
type Config struct {
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
		SSLMode  string
	}
}

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt возвращает целочисленное значение переменной окружения или значение по умолчанию
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// loadConfig загружает конфигурацию из переменных окружения
func loadConfig() *Config {
	var cfg Config

	// Database configuration
	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getEnvInt("DB_PORT", 5432)
	cfg.Database.User = getEnv("DB_USER", "postgres")
	cfg.Database.Password = getEnv("DB_PASSWORD", "postgres")
	cfg.Database.Name = getEnv("DB_NAME", "user_service")
	cfg.Database.SSLMode = getEnv("DB_SSLMODE", "disable")

	return &cfg
}

// buildDSN строит строку подключения к PostgreSQL
func buildDSN(cfg *Config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)
}

func main() {
	// Загружаем конфигурацию из переменных окружения
	cfg := loadConfig()

	log.Printf("Using database: %s@%s:%d/%s",
		cfg.Database.User,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name)

	// Строим DSN
	dsn := buildDSN(cfg)

	// Подключаемся к базе данных
	db, err := sql.Open("pg", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Создаем таблицу для отслеживания миграций
	if err := createMigrationsTable(db); err != nil {
		log.Fatalf("Failed to create migrations table: %v", err)
	}

	// Применяем миграции
	if err := runMigrations(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("✅ All migrations applied successfully")
}

func createMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP DEFAULT NOW(),
			UNIQUE(name)
		)
	`
	_, err := db.Exec(query)
	return err
}

func runMigrations(db *sql.DB) error {
	migrationsDir := "db/migrations"

	// Проверяем существование директории
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		log.Printf("⚠️  Migrations directory '%s' does not exist, creating...", migrationsDir)
		if err := os.MkdirAll(migrationsDir, 0755); err != nil {
			return fmt.Errorf("failed to create migrations directory: %v", err)
		}
		log.Printf("✅ Created migrations directory: %s", migrationsDir)
		return nil
	}

	// Читаем файлы миграций
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %v", err)
	}

	// Фильтруем и сортируем SQL файлы
	var migrationFiles []fs.DirEntry
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file)
		}
	}

	if len(migrationFiles) == 0 {
		log.Println("ℹ️  No migration files found")
		return nil
	}

	// Сортируем по возрастанию (001_, 002_, ...)
	sort.Slice(migrationFiles, func(i, j int) bool {
		return migrationFiles[i].Name() < migrationFiles[j].Name()
	})

	// Применяем миграции
	for _, file := range migrationFiles {
		migrationName := file.Name()

		// Проверяем, применена ли уже миграция
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM migrations WHERE name = $1", migrationName).Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to check migration status: %v", err)
		}
		if count > 0 {
			log.Printf("⏭️  Migration %s already applied, skipping", migrationName)
			continue
		}

		// Читаем SQL файл
		filePath := filepath.Join(migrationsDir, migrationName)
		sqlContent, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %v", migrationName, err)
		}

		// Выполняем миграцию в транзакции
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %v", err)
		}

		log.Printf("🔄 Applying migration: %s", migrationName)

		if _, err := tx.Exec(string(sqlContent)); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %v", migrationName, err)
		}

		// Записываем в таблицу миграций
		_, err = tx.Exec("INSERT INTO migrations (name, applied_at) VALUES ($1, $2)",
			migrationName, time.Now())
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update migrations table: %v", err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction: %v", err)
		}

		log.Printf("✅ Migration %s applied successfully", migrationName)
	}

	return nil
}
