package handlers

import (
	"math"
	"fmt"
	"encoding/csv"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"payslip-system/config"
	"payslip-system/models"
	"payslip-system/tax"
)

// Create payslip handler
func CreatePayslip(c *gin.Context) {
	name := c.PostForm("employee_name")
	salaryStr := c.PostForm("annual_salary")
	salary, _ := strconv.Atoi(salaryStr)

	monthlyTax := tax.GenerateMonthlyPayslip(salary, 0) / 12.0
	grossMonthly := float64(salary) / 12
	netMonthly := grossMonthly - monthlyTax

	payslip := models.Payslip{
		EmployeeName:       name,
		AnnualSalary:       salary,
		GrossMonthlyIncome: grossMonthly,
		MonthlyIncomeTax:   monthlyTax,
		NetMonthlyIncome:   netMonthly,
		CreatedAt:          time.Now(),
	}

	config.DB.Create(&payslip)
	c.Redirect(http.StatusFound, fmt.Sprintf("/payslip/%d", payslip.ID))
}

// Get all payslips data
func GetAllPayslips(c *gin.Context) {
	var payslips []models.Payslip
	config.DB.Find(&payslips)
	c.JSON(http.StatusOK, payslips)
}

// Get data By ID
func GetPayslipByID(c *gin.Context) {
    id := c.Param("id")

    var payslip models.Payslip
    if err := config.DB.First(&payslip, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Payslip not found"})
        return
    }

    response := fmt.Sprintf(
		"Monthly Payslip for: %s\nGross Monthly Income: $%.2f\nMonthly Income Tax: $%.2f\nNet Monthly Income: $%.2f",
		payslip.EmployeeName,
		payslip.GrossMonthlyIncome,
		payslip.MonthlyIncomeTax,
		payslip.NetMonthlyIncome,
	)
	

    c.String(http.StatusOK, response)
}

// Delete By ID
func DeletePayslip(c *gin.Context) {
    id := c.Param("id")
    config.DB.Delete(&models.Payslip{}, id)
    c.JSON(200, gin.H{"deleted": id})
}

// Generate CSV File
func ExportCSV(c *gin.Context) {
    var payslips []models.Payslip
    config.DB.Find(&payslips)

    c.Header("Content-Type", "text/csv")
    // Set the name of file CSV
    c.Header("Content-Disposition", "attachment; filename=monthly_payslip.csv")

    writer := csv.NewWriter(c.Writer)
    defer writer.Flush()

    writer.Write([]string{"ID", "Name", "Annual Salary", "Gross Monthly Salary", "Monthly Tax", "Net Monthly Salary"})

    for _, p := range payslips {
        writer.Write([]string{
            strconv.Itoa(int(p.ID)),
            p.EmployeeName,
            strconv.Itoa(p.AnnualSalary),
            strconv.FormatFloat(float64(p.GrossMonthlyIncome), 'f', 2, 64),
            strconv.FormatFloat(float64(p.MonthlyIncomeTax), 'f', 2, 64),
            strconv.FormatFloat(float64(p.NetMonthlyIncome), 'f', 2, 64),
        })
    }
}

// Function to make .00 format
func Round(f float64, places int) float64 {
    shift := math.Pow(10, float64(places))
    return math.Round(f*shift) / shift
}