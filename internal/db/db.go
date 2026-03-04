package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	paths := []string{".env", "../.env", "../../.env"}
	for _, path := range paths {
		if err := godotenv.Load(path); err == nil {
			fmt.Println("Loaded .env from", path)
			return
		}
	}
	log.Println("Warning: .env file not found in any location")
}

func ConnectPostgres(ctx context.Context) (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	if sslmode == "" {
		sslmode = "require"
	}

	if host == "" || port == "" || user == "" || password == "" || dbName == "" {
		return nil, fmt.Errorf("missing required database environment variables")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbName, sslmode)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database ping successful")
	log.Println("Successfully connected to database")
	return DB, nil
}

func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return fmt.Errorf("failed to get sql.DB: %w", err)
		}
		if err := sqlDB.Close(); err != nil {
			return fmt.Errorf("failed to close database: %w", err)
		}
		log.Println("Database connection closed")
	}
	return nil
}