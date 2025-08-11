package dto

type ReportFilterRequest struct {
	StartDate  string `query:"start_date"`
	EndDate    string `query:"end_date"`
	Filter     string `query:"filter"`
	Status     string `query:"status"`
	MerchantID string `json:"-"`
	UserID     string `json:"-"`
}

type DailyReportResponse struct {
	Label            string        `json:"label"`      // e.g., "2025-07-08" or "2025-07-08 to 2025-07-14"
	StartDate        string        `json:"start_date"` // optional, mostly for weekly
	EndDate          string        `json:"end_date"`   // optional, mostly for weekly
	TotalTransaction int           `json:"totalTransaction"`
	Paid             float64       `json:"paid"`
	Unpaid           float64       `json:"unpaid"`
	TotalIncome      float64       `json:"totalIncome"`
	MerchantID       string        `json:"merchant_id"`
	BestSellers      []BestSeller  `json:"best_sellers"` // ← Tambahkan
	Products         []ProductItem `json:"products"`
}

type BestSeller struct {
	ProductName string `json:"product_name"`
	TotalQty    int    `json:"total_qty"`
}

type ProductItem struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}
