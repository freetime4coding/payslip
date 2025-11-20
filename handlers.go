package main

import (
    "fmt"
    "math"
    "time"
    "encoding/csv"
    "github.com/gin-gonic/gin"
    "strconv"
    "html/template"
    "net/http"
)

var tmpl = template.Must(template.ParseGlob("templates/*.html"))

// Form handler
func FormHandler(w http.ResponseWriter, r *http.Request) {
    tmpl.ExecuteTemplate(w, "form.html", nil)
}

// List handler
func ListHandler(w http.ResponseWriter, r *http.Request) {
    tmpl.ExecuteTemplate(w, "list.html", nil)
}

// Create payslip function
func CreatePayslip(c *gin.Context) {
	var input struct {
		EmployeeName string `form:"employee_name" binding:"required"`
		AnnualSalary int    `form:"annual_salary" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.HTML(400, "form.html", gin.H{"error": err.Error()})
		return
	}

	annualTax := GenerateMonthlyPayslip(input.AnnualSalary, 0)
	grossMonthly := float64(input.AnnualSalary) / 12
	monthlyTax := annualTax / 12
	netMonthly := grossMonthly - monthlyTax

	payslip := Payslip{
		EmployeeName:           input.EmployeeName,
		AnnualSalary:           input.AnnualSalary,
        GrossMonthlyIncome:     Round(grossMonthly, 2),
        MonthlyIncomeTax:       Round(monthlyTax, 2),
        NetMonthlyIncome:       Round(netMonthly, 2),
        CreatedAt:              time.Now(),
	}

	DB.Create(&payslip)

    payslipIDStr := strconv.Itoa(int(payslip.ID))
    redirectURL := "/payslip/" + payslipIDStr
	c.Redirect(302, redirectURL)
}

// Get all data
func GetAllPayslips(c *gin.Context) {
    var payslips []Payslip
    DB.Find(&payslips)

    c.JSON(200, payslips)
}

// Get data By ID
func GetPayslipByID(c *gin.Context) {
    id := c.Param("id")

    var payslip Payslip
    if err := DB.First(&payslip, id).Error; err != nil {
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
    DB.Delete(&Payslip{}, id)
    c.JSON(200, gin.H{"deleted": id})
}

// Generate CSV File
func ExportCSV(c *gin.Context) {
    var payslips []Payslip
    DB.Find(&payslips)

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