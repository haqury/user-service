package main

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/haqury/user-service/internal/config"
	_ "github.com/uptrace/bun/driver/pgdriver" // PostgreSQL драйвер
)

func main() {
	// Создаем конфигурацию
	c, err := config.NewConfig("", "")
	if err != nil {
		log.Fatalf("Failed to create config: %v", err)
	}

	log.Printf("Using database: %s", c.Database.DSN)

	// Подключаемся к базе данных
	db, err := sql.Open("pg", c.Database.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Создаем таблицу для отслеживания seeds
	if err := createSeedsTable(db); err != nil {
		log.Fatalf("Failed to create seeds table: %v", err)
	}

	// Применяем seeds
	if err := runSeeds(db); err != nil {
		log.Fatalf("Seeds failed: %v", err)
	}

	log.Println("✅ All seeds applied successfully")
}

func createSeedsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS seeds (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP DEFAULT NOW(),
			UNIQUE(name)
		)
	`
	_, err := db.Exec(query)
	return err
}

func runSeeds(db *sql.DB) error {
	seedsDir := "db/seeds"

	// Проверяем существование директории
	if _, err := os.Stat(seedsDir); os.IsNotExist(err) {
		log.Printf("⚠️  Seeds directory '%s' does not exist, creating...", seedsDir)
		if err := os.MkdirAll(seedsDir, 0755); err != nil {
			return fmt.Errorf("failed to create seeds directory: %v", err)
		}
		log.Printf("✅ Created seeds directory: %s", seedsDir)
		return nil
	}

	// Читаем файлы seeds
	files, err := os.ReadDir(seedsDir)
	if err != nil {
		return fmt.Errorf("failed to read seeds directory: %v", err)
	}

	// Фильтруем и сортируем SQL файлы
	var seedFiles []fs.DirEntry
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			seedFiles = append(seedFiles, file)
		}
	}

	if len(seedFiles) == 0 {
		log.Println("ℹ️  No seed files found")
		return nil
	}

	// Сортируем по возрастанию (001_, 002_, ...)
	sort.Slice(seedFiles, func(i, j int) bool {
		return seedFiles[i].Name() < seedFiles[j].Name()
	})

	// Применяем seeds
	for _, file := range seedFiles {
		seedName := file.Name()

		// Проверяем, применен ли уже seed
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM seeds WHERE name = $1", seedName).Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to check seed status: %v", err)
		}
		if count > 0 {
			log.Printf("⏭️  Seed %s already applied, skipping", seedName)
			continue
		}

		// Читаем SQL файл
		filePath := filepath.Join(seedsDir, seedName)
		sqlContent, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read seed file %s: %v", seedName, err)
		}

		// Выполняем seed в транзакции
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %v", err)
		}

		log.Printf("🌱 Applying seed: %s", seedName)

		if _, err := tx.Exec(string(sqlContent)); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute seed %s: %v", seedName, err)
		}

		// Записываем в таблицу seeds
		_, err = tx.Exec("INSERT INTO seeds (name, applied_at) VALUES ($1, $2)",
			seedName, time.Now())
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update seeds table: %v", err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction: %v", err)
		}

		log.Printf("✅ Seed %s applied successfully", seedName)
	}

	return nil
}
