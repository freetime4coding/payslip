package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"payslip-system/config"
	"payslip-system/models"
	"payslip-system/handlers"
)

func setupTestDB(t *testing.T) *gorm.DB {
	os.Setenv("APP_ENV", "test")
	config.LoadConfig()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.C.DBHost, config.C.DBUser, config.C.DBPassword, config.C.DBName, config.C.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect PostgreSQL for tests: %v", err)
	}

	db.Migrator().DropTable(&models.Payslip{})
	db.AutoMigrate(&models.Payslip{})
	config.DB = db
	return db
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/payslip", handlers.CreatePayslip)
	r.GET("/payslip", handlers.GetAllPayslips)
	return r
}

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
	config.DB.Model(&models.Payslip{}).Count(&count)
	if count != 1 {
		t.Fatalf("expected 1 payslip saved, got %d", count)
	}
}

func TestGetAllPayslips(t *testing.T) {
	setupTestDB(t)
	config.DB.Create(&models.Payslip{
		EmployeeName:       "Tester",
		AnnualSalary:       120000,
		GrossMonthlyIncome: 10000,
		MonthlyIncomeTax:   1000,
		NetMonthlyIncome:   9000,
		CreatedAt:          time.Now(),
	})

	router := setupRouter()
	req, _ := http.NewRequest("GET", "/payslip", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}

	var result []models.Payslip
	err := json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 payslip, got %d", len(result))
	}
}
