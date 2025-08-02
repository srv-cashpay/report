package dto

type ReportFilterRequest struct {
	StartDate string `query:"start_date"`
	EndDate   string `query:"end_date"`
	Filter    string `query:"filter"`
	Status    string `query:"status"` // tambah ini
}

type DailyReportResponse struct {
	Date             string `json:"date"`
	TotalTransaction int    `json:"totalTransaction"`
	Paid             int    `json:"paid"`
	Unpaid           int    `json:"unpaid"`
	TotalIncome      int    `json:"totalIncome"`
}
