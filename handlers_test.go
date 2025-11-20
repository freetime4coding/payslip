package main

import (
    "testing"
    "net/http"
    "net/http/httptest"
    "github.com/gin-gonic/gin"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "strings"
    "encoding/json"
    "time"
)

// Init test DB
func setupTestDB(t *testing.T) *gorm.DB {
    dsn := "host=localhost user=postgres password=YOUR_PASSWORD dbname=payslipstestdb port=5432 sslmode=disable"

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        t.Fatalf("failed to connect PostgreSQL for tests: %v", err)
    }

    db.Migrator().DropTable(&Payslip{})
    db.AutoMigrate(&Payslip{})

    DB = db
    return db
}

// Set list of Test API
func setupRouter() *gin.Engine {
    gin.SetMode(gin.TestMode)
    r := gin.Default()

    r.POST("/payslip", CreatePayslip)
    r.GET("/payslip", GetAllPayslips)

    return r
}

// Run the test create payslip
func TestCreatePayslip(t *testing.T) {
    setupTestDB(t)
    router := setupRouter()

    form := strings.NewReader("employee_name=Test&annual_salary=60000")

    req, _ := http.NewRequest("POST", "/payslip", form)
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    if w.Code != http.StatusFound {
        t.Fatalf("expected 302 redirect, got %d", w.Code)
    }

    var count int64
    DB.Model(&Payslip{}).Count(&count)
    if count != 1 {
        t.Fatalf("expected 1 payslip saved, got %d", count)
    }
}

// Run the test get list payslip
func TestGetAllPayslips(t *testing.T) {
    setupTestDB(t)

    DB.Create(&Payslip{
        EmployeeName:        "Tester",
        AnnualSalary:        100000,
        GrossMonthlyIncome:  8333.33,
        MonthlyIncomeTax:    833.33,
        NetMonthlyIncome:    7500,
        CreatedAt:           time.Now(),
    })

    router := setupRouter()

    req, _ := http.NewRequest("GET", "/payslip", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Fatalf("expected 200 OK, got %d", w.Code)
    }

    var result []Payslip
    if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
        t.Fatalf("invalid json: %v", err)
    }

    if len(result) != 1 {
        t.Fatalf("expected 1 payslip, got %d", len(result))
    }
}