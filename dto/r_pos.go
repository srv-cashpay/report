package dto

type ReportFilterRequest struct {
	StartDate string `query:"start_date"`
	EndDate   string `query:"end_date"`
	Filter    string `query:"filter"`
}

type DailyReportResponse struct {
	Date             string `json:"date"`
	TotalTransaction string `json:"total_transaction"`
	TotalIncome      string `json:"total_income"`
}
