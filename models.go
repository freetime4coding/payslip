package main

import (
    "time"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var DB *gorm.DB

// Initialize the Object
type Payslip struct {
    ID                    uint      `json:"id" gorm:"primaryKey"`
    EmployeeName          string    `json:"employee_name"`
    AnnualSalary          int       `json:"annual_salary,string"`
    GrossMonthlyIncome    float64   `json:"gross_monthly_income,string"`
    NetMonthlyIncome      float64   `json:"net_monthly_income,string"`
    MonthlyIncomeTax      float64   `json:"monthly_income_tax,string"`
    CreatedAt             time.Time `json:"created_at"`
}

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
