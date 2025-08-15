package product

import (
	"fmt"
	"time"

	"github.com/srv-cashpay/report/dto"
)

func (r *rproductReproductitory) SlowMoving(start, end time.Time, filter, status, merchantID string) ([]dto.BestSellerResponse, error) {
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
			COALESCE(SUM(CASE WHEN status_payment = 'Unpaid' THEN pay ELSE 0 END), 0) AS unpaid,
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

	var bestSellers []dto.BestSeller
	if err := r.DB.Raw(`
	SELECT 
	(json_data->>'product_name') AS product_name,
	SUM((json_data->>'quantity')::int) AS total_qty
	FROM pos,
	LATERAL json_array_elements(product::json) AS json_data
	WHERE merchant_id = ?
	AND status_payment = 'Paid'
	AND created_at AT TIME ZONE 'Asia/Jakarta' >= ?
	AND created_at AT TIME ZONE 'Asia/Jakarta' <= ?
	AND deleted_at IS NULL
	GROUP BY product_name
	ORDER BY total_qty ASC
	LIMIT 10
`, merchantID, startDate, endDate).Scan(&bestSellers).Error; err != nil {
		return nil, err
	}
	var result []dto.BestSellerResponse
	for _, row := range rows {
		var label, startDateStr, endDateStr string

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

		result = append(result, dto.BestSellerResponse{
			Label:            label,
			StartDate:        startDateStr,
			EndDate:          endDateStr,
			TotalTransaction: row.TotalTransaction,
			Paid:             row.Paid,
			Unpaid:           row.Unpaid,
			TotalIncome:      row.TotalIncome,
			MerchantID:       row.MerchantID,
			BestSellers:      bestSellers, // ← tambahkan ini di DTO
		})
	}

	return result, nil
}
