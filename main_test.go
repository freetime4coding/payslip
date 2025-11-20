package main

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

// Init test DB
func setupTestDBForMain(t *testing.T) {
    dsn := "host=localhost user=postgres password=YOUR_PASSWORD dbname=payslipstestdb port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        t.Fatalf("open db: %v", err)
    }
    DB = db
    DB.AutoMigrate(&Payslip{})
}

// Set the list of Test API
func setupMainRouter() *gin.Engine {
    gin.SetMode(gin.TestMode)
    r := gin.Default()

    r.POST("/payslip", CreatePayslip)
    r.GET("/payslip", GetAllPayslips)
    r.GET("/payslip/:id", GetPayslipByID)
    r.DELETE("/payslip/:id", DeletePayslip)
    r.GET("/payslip/export", ExportCSV)

    return r
}

// Run the test route and function
func TestRoutesWiring(t *testing.T) {
    setupTestDBForMain(t)
    router := setupMainRouter()

    // Test GET /payslip (empty list)
    req, _ := http.NewRequest("GET", "/payslip", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d", w.Code)
    }
}
