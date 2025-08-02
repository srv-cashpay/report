package pos

import (
	"time"

	"github.com/srv-cashpay/report/dto"
)

func (r *rposRepository) GetReport(start, end time.Time) ([]dto.DailyReportResponse, error) {
	var reports []dto.DailyReportResponse

	err := r.DB.
		Raw(`
			SELECT 
				DATE(created_at) as date, 
				COUNT(*) as total_transaction, 
				SUM(amount) as total_income
			FROM transactions
			WHERE DATE(created_at) BETWEEN ? AND ?
			GROUP BY DATE(created_at)
			ORDER BY DATE(created_at)
		`, start.Format("2006-01-02"), end.Format("2006-01-02")).
		Scan(&reports).Error

	if err != nil {
		return nil, err
	}

	return reports, nil
}
