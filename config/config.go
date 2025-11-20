package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config struct holds environment variables
type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
}

// Global variables
var C *Config
var DB *gorm.DB

// LoadConfig loads DB configuration from environment variables
func LoadConfig() {
	C = &Config{
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
	}

	if C.DBUser == "" || C.DBPassword == "" || C.DBName == "" || C.DBHost == "" || C.DBPort == "" {
		log.Fatal("Missing database configuration in environment")
	}
}

// InitDB initializes the global database connection
func InitDB() {
	if C == nil {
		LoadConfig()
	}

	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		C.DBUser, C.DBPassword, C.DBName, C.DBHost, C.DBPort,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("failed to get database object: %v", err)
	}

	// Optional connection pool configuration
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("[info] database connected successfully")
}