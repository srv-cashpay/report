package pos

import (
	"fmt"
	"time"

	"github.com/srv-cashpay/report/dto"
)

func (r *rposRepository) GetReport(start, end time.Time, filter, status, merchantID string) ([]dto.DailyReportResponse, error) {
	var rows []struct {
		GroupDate        time.Time
		TotalTransaction int
		Paid             float64
		Unpaid           float64
		TotalIncome      float64
	}

	var groupExpr string
	switch filter {
	case "weekly":
		groupExpr = "DATE_TRUNC('week', created_at)"
	case "monthly":
		groupExpr = "DATE_TRUNC('month', created_at)"
	case "yearly":
		groupExpr = "DATE_TRUNC('year', created_at)"
	default:
		groupExpr = "DATE_TRUNC('day', created_at)"
	}

	query := fmt.Sprintf(`
		SELECT 
			%s AS group_date,
			COUNT(*) AS total_transaction,
			COALESCE(SUM(CASE WHEN status_payment = 'Paid' THEN pay ELSE 0 END), 0) AS paid,
			COALESCE(SUM(CASE WHEN status_payment = 'Unpaid' THEN pay ELSE 0 END), 0) AS unpaid,
			COALESCE(SUM(pay), 0) AS total_income
		FROM pos
		WHERE created_at BETWEEN ? AND ?
			AND deleted_at IS NULL
					AND merchant_id = ?
	`, groupExpr)

	args := []interface{}{start, end, merchantID}

	if status != "" {
		query += " AND status_payment = ?"
		args = append(args, status)
	}

	query += " GROUP BY group_date ORDER BY group_date"

	if err := r.DB.Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}

	var result []dto.DailyReportResponse
	for _, row := range rows {
		var label, startDate, endDate string

		switch filter {
		case "weekly":
			start := row.GroupDate
			end := start.AddDate(0, 0, 6)
			label = fmt.Sprintf("%s to %s", start.Format("2006-01-02"), end.Format("2006-01-02"))
			startDate = start.Format("2006-01-02")
			endDate = end.Format("2006-01-02")
		case "monthly":
			year, month, _ := row.GroupDate.Date()
			start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
			end := start.AddDate(0, 1, -1)
			label = fmt.Sprintf("%s to %s", start.Format("2006-01-02"), end.Format("2006-01-02"))
			startDate = start.Format("2006-01-02")
			endDate = end.Format("2006-01-02")
		case "yearly":
			year := row.GroupDate.Year()
			start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
			end := time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC)
			label = fmt.Sprintf("%d", year)
			startDate = start.Format("2006-01-02")
			endDate = end.Format("2006-01-02")
		default: // daily
			label = row.GroupDate.Format("2006-01-02")
			startDate = label
			endDate = label
		}

		result = append(result, dto.DailyReportResponse{
			Label:            label,
			StartDate:        startDate,
			EndDate:          endDate,
			TotalTransaction: row.TotalTransaction,
			Paid:             row.Paid,
			Unpaid:           row.Unpaid,
			TotalIncome:      row.TotalIncome,
		})
	}

	return result, nil
}
