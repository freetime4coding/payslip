package main

import (
	"net/http"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    r.LoadHTMLGlob("templates/*.html")
	
	// Route: View home page
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "form.html", gin.H{
			"title": "Create Payslip",
		})
	})
	
	// Route: View list of payslip page
	r.GET("/list", func(c *gin.Context) {
		var payslips []Payslip
		DB.Order("id desc").Find(&payslips)
		c.HTML(200, "list.html", gin.H{
			"payslips": payslips,
		})
	})

	// Route API: Create Payslip (submit data)
    r.POST("/payslip", CreatePayslip)

	// Route API: get all data payslip
	r.GET("/payslip", GetAllPayslips)

	// Route API: get data By ID
	r.GET("/payslip/:id", GetPayslipByID)
	
	// Route API: delete data By ID
	r.DELETE("/payslip/:id", DeletePayslip)

	// Route API: generate csv list of payslips 
    r.GET("/payslip/export", ExportCSV)
    
	// set the port to run the system
	r.Run(":3000")
}