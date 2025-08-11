package pos

import (
	"fmt"
	"time"

	"github.com/srv-cashpay/report/dto"
)

func (r *rposRepository) Order(start, end time.Time, filter, status, merchantID string) ([]dto.DailyReportResponse, error) {
	var rows []struct {
		GroupDate        time.Time
		TotalTransaction int
		Paid             float64
		Unpaid           float64
		TotalIncome      float64
		MerchantID       string
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	startDate := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, loc)
	endDate := time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, loc)

	var groupExpr string
	switch filter {
	case "weekly":
		groupExpr = "DATE_TRUNC('week', created_at AT TIME ZONE 'Asia/Jakarta')"
	case "monthly":
		groupExpr = "DATE_TRUNC('month', created_at AT TIME ZONE 'Asia/Jakarta')"
	case "yearly":
		groupExpr = "DATE_TRUNC('year', created_at AT TIME ZONE 'Asia/Jakarta')"
	default:
		groupExpr = "DATE_TRUNC('day', created_at AT TIME ZONE 'Asia/Jakarta')"
	}

	query := fmt.Sprintf(`
		SELECT 
			%s AS group_date,
			COUNT(*) AS total_transaction,
			COALESCE(SUM(CASE WHEN status_payment = 'Paid' THEN pay ELSE 0 END), 0) AS paid,
			COALESCE(SUM(
				CASE 
					WHEN status_payment = 'Unpaid' THEN (
						SELECT SUM(
							(COALESCE((elem->>'price')::numeric, 0)) * 
							(COALESCE((elem->>'quantity')::int, 0))
						)
						FROM jsonb_array_elements(product::jsonb) AS elem
					)
					ELSE 0
				END
			), 0) AS unpaid,
			COALESCE(SUM(pay), 0) AS total_income,
			merchant_id
		FROM pos
		WHERE created_at AT TIME ZONE 'Asia/Jakarta' >= ? 
			AND created_at AT TIME ZONE 'Asia/Jakarta' <= ?
			AND deleted_at IS NULL
			AND merchant_id = ?
	`, groupExpr)

	args := []interface{}{startDate, endDate, merchantID}
	if status != "" {
		query += " AND status_payment = ?"
		args = append(args, status)
	}

	query += " GROUP BY group_date, merchant_id ORDER BY group_date"

	if err := r.DB.Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}

	var result []dto.DailyReportResponse
	for _, row := range rows {
		var label, startDateStr, endDateStr string

		// Tentukan range tanggal per group
		switch filter {
		case "weekly":
			start := row.GroupDate.In(loc)
			end := start.AddDate(0, 0, 6)
			label = fmt.Sprintf("%s to %s", start.Format("2006-01-02"), end.Format("2006-01-02"))
			startDateStr = start.Format("2006-01-02")
			endDateStr = end.Format("2006-01-02")
		case "monthly":
			year, month, _ := row.GroupDate.In(loc).Date()
			start := time.Date(year, month, 1, 0, 0, 0, 0, loc)
			end := start.AddDate(0, 1, -1)
			label = fmt.Sprintf("%s to %s", start.Format("2006-01-02"), end.Format("2006-01-02"))
			startDateStr = start.Format("2006-01-02")
			endDateStr = end.Format("2006-01-02")
		case "yearly":
			year := row.GroupDate.In(loc).Year()
			start := time.Date(year, 1, 1, 0, 0, 0, 0, loc)
			end := time.Date(year, 12, 31, 0, 0, 0, 0, loc)
			label = fmt.Sprintf("%d", year)
			startDateStr = start.Format("2006-01-02")
			endDateStr = end.Format("2006-01-02")
		default:
			groupDate := row.GroupDate.In(loc)
			label = groupDate.Format("2006-01-02")
			startDateStr = label
			endDateStr = label
		}

		// Ambil semua produk pada periode ini
		var products []dto.ProductItem
		productQuery := `
			SELECT 
				elem->>'product_id' AS product_id,
				elem->>'product_name' AS product_name,
				(elem->>'quantity')::int AS quantity,
				(elem->>'price')::numeric AS price
			FROM pos,
			jsonb_array_elements(product::jsonb) AS elem
			WHERE merchant_id = ?
			AND created_at AT TIME ZONE 'Asia/Jakarta' >= ? 
			AND created_at AT TIME ZONE 'Asia/Jakarta' <= ?
			AND deleted_at IS NULL
		`
		if err := r.DB.Raw(productQuery, row.MerchantID, startDateStr, endDateStr).Scan(&products).Error; err != nil {
			return nil, err
		}

		// Gabungkan ke hasil akhir
		result = append(result, dto.DailyReportResponse{
			Label:            label,
			StartDate:        startDateStr,
			EndDate:          endDateStr,
			TotalTransaction: row.TotalTransaction,
			Paid:             row.Paid,
			Unpaid:           row.Unpaid,
			TotalIncome:      row.TotalIncome,
			MerchantID:       row.MerchantID,
			Products:         products,
		})
	}

	return result, nil
}
