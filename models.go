package main

import (
    "time"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var DB *gorm.DB

// Set DB
func init() {
    dsn := "host=localhost user=postgres password=YOUR_PASSWORD dbname=payslipdb port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Failed to connect to database")
    }
    DB = db
    DB.AutoMigrate(&Payslip{})
}
