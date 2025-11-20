package models

import "time"

type Payslip struct {
	ID                 uint      `json:"id" gorm:"primaryKey"`
	EmployeeName       string    `json:"employee_name"`
	AnnualSalary       int       `json:"annual_salary,string"`
	GrossMonthlyIncome float64   `json:"gross_monthly_income,string"`
	NetMonthlyIncome   float64   `json:"net_monthly_income,string"`
	MonthlyIncomeTax   float64   `json:"monthly_income_tax,string"`
	CreatedAt          time.Time `json:"created_at"`
}
