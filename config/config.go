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

// Config struct to hold environment configuration
type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
}

// Global config variable
var C *Config

// Global DB connection
var DB *gorm.DB

// LoadConfig loads environment variables into C
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

// InitDB initializes the database connection
func InitDB() {
	// Load config if not already loaded
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
		log.Fatalf("[error] failed to initialize database, got error %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("[error] failed to get database object, got error %v", err)
	}

	// Optional: configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("[info] database connected successfully")
}
