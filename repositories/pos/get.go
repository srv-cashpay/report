package pos

import (
	"time"

	"github.com/srv-cashpay/report/dto"
)

func (r *rposRepository) GetReport(start, end time.Time, filter, status string) ([]dto.DailyReportResponse, error) {
	var reports []dto.DailyReportResponse

	var groupBy string
	switch filter {
	case "weekly":
		groupBy = "TO_CHAR(created_at, 'IYYY-IW')"
	case "monthly":
		groupBy = "TO_CHAR(created_at, 'YYYY-MM')"
	case "yearly":
		groupBy = "TO_CHAR(created_at, 'YYYY')"
	default:
		groupBy = "DATE(created_at)"
	}

	query := `
		SELECT 
			` + groupBy + ` AS date,
			COUNT(*) AS total_transaction,
			SUM(pay) AS total_income
		FROM pos
		WHERE DATE(created_at) BETWEEN ? AND ?
			AND deleted_at IS NULL
	`

	args := []interface{}{start.Format("2006-01-02"), end.Format("2006-01-02")}

	if status != "" {
		query += " AND status_payment = ?"
		args = append(args, status)
	}

	query += " GROUP BY " + groupBy + " ORDER BY " + groupBy

	err := r.DB.Raw(query, args...).Scan(&reports).Error
	if err != nil {
		return nil, err
	}

	return reports, nil
}
