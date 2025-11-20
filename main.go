package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"payslip-system/models"
	"payslip-system/config"
	"payslip-system/handlers"
)

func main() {
	config.InitDB()

	r := gin.Default()
    r.LoadHTMLGlob("templates/*.html")
	
	// Home page
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "form.html", gin.H{
			"title": "Create Payslip",
		})
	})

	// List page
	r.GET("/list", func(c *gin.Context) {
		var payslips []models.Payslip
		config.DB.Order("id desc").Find(&payslips)
		c.HTML(http.StatusOK, "list.html", gin.H{
			"payslips": payslips,
		})
	})
	
	r.POST("/payslip", handlers.CreatePayslip) 			// Route API: Create Payslip (submit data)
	r.GET("/payslip/:id", handlers.GetPayslipByID)		// Route API: get all data payslip
	r.GET("/payslips", handlers.GetAllPayslips) 		// Route API: get data By ID
	r.DELETE("/payslip/:id", handlers.DeletePayslip)	// Route API: delete data By ID
	r.GET("/export-csv", handlers.ExportCSV) 			// Route API: generate csv list of payslips

	r.Run(":3000") // set the port to run the system
}