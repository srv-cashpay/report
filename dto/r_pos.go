package dto

type ReportFilterRequest struct {
	StartDate string `query:"start_date"`
	EndDate   string `query:"end_date"`
	Filter    string `query:"filter"`
}

type DailyReportResponse struct {
	Date             string `json:"date" gorm:"column:date"`
	TotalTransaction int    `json:"totalTransaction" gorm:"column:total_transaction"`
	TotalIncome      int64  `json:"totalIncome" gorm:"column:total_income"`
	Paid             int    `json:"paid"`
	Unpaid           int    `json:"unpaid"`
}
