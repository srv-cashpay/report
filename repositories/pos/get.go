package pos

import (
	"fmt"
	"time"

	"github.com/srv-cashpay/report/dto"
)

func (r *rposRepository) GetReport(start, end time.Time, rangeType string) ([]dto.DailyReportResponse, error) {
	var reports []dto.DailyReportResponse

	groupBy := ""
	dateSelect := ""

	switch rangeType {
	case "daily":
		groupBy = "DATE(created_at)"
		dateSelect = "DATE(created_at) AS date"
	case "weekly":
		groupBy = "YEAR(created_at), WEEK(created_at, 1)"
		dateSelect = "CONCAT(YEAR(created_at), '-W', LPAD(WEEK(created_at, 1), 2, '0')) AS date"
	case "monthly":
		groupBy = "YEAR(created_at), MONTH(created_at)"
		dateSelect = "DATE_FORMAT(created_at, '%Y-%m') AS date"
	case "yearly":
		groupBy = "YEAR(created_at)"
		dateSelect = "YEAR(created_at) AS date"
	default:
		return nil, fmt.Errorf("invalid range type")
	}

	query := fmt.Sprintf(`
		SELECT
			%s,
			COUNT(*) AS total_transaction,
			SUM(CASE WHEN status_payment = 'Paid' THEN pay ELSE 0 END) AS total_income,
			COUNT(CASE WHEN status_payment = 'Paid' THEN 1 END) AS paid,
			COUNT(CASE WHEN status_payment = 'Unpaid' THEN 1 END) AS unpaid
		FROM pos
		WHERE DATE(created_at) BETWEEN ? AND ?
			AND deleted_at IS NULL
		GROUP BY %s
		ORDER BY date
	`, dateSelect, groupBy)

	err := r.DB.Raw(query, start.Format("2006-01-02"), end.Format("2006-01-02")).Scan(&reports).Error
	if err != nil {
		return nil, err
	}

	return reports, nil
}
