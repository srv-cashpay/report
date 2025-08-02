package pos

import (
	"fmt"
	"time"

	"github.com/srv-cashpay/report/dto"
)

func (r *rposRepository) GetReport(start, end time.Time, filter, status string) ([]dto.DailyReportResponse, error) {
	var reports []dto.DailyReportResponse

	// Tentukan grouping berdasarkan filter
	var groupBy string
	switch filter {
	case "weekly":
		groupBy = "TO_CHAR(created_at, 'IYYY-IW')" // ISO week: e.g., "2025-31"
	case "monthly":
		groupBy = "TO_CHAR(created_at, 'YYYY-MM')" // e.g., "2025-08"
	case "yearly":
		groupBy = "EXTRACT(YEAR FROM created_at)::text"
	default:
		groupBy = "TO_CHAR(created_at, 'YYYY-MM-DD')" // daily default
	}

	// SQL query dinamis
	query := fmt.Sprintf(`
		SELECT 
			%s AS date,
			COUNT(*) AS total_transaction,
			COALESCE(SUM(CASE WHEN status_payment = 'Paid' THEN pay ELSE 0 END), 0) AS paid,
			COALESCE(SUM(CASE WHEN status_payment = 'Unpaid' THEN pay ELSE 0 END), 0) AS unpaid,
			COALESCE(SUM(pay), 0) AS total_income
		FROM pos
		WHERE created_at BETWEEN ? AND ?
			AND deleted_at IS NULL
	`, groupBy)

	// Tambahkan filter status jika ada
	var args []interface{}
	args = append(args, start, end)

	if status != "" {
		query += " AND status_payment = ?"
		args = append(args, status)
	}

	query += fmt.Sprintf(" GROUP BY %s ORDER BY %s", groupBy, groupBy)

	// Eksekusi query
	if err := r.DB.Raw(query, args...).Scan(&reports).Error; err != nil {
		return nil, err
	}

	return reports, nil
}
