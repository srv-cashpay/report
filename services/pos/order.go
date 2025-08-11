package pos

import (
	"time"

	"github.com/srv-cashpay/report/dto"
)

func (s *rposService) Order(start, end time.Time, filter, status, merchantID string) ([]dto.DailyReportResponse, error) {
	return s.Repo.Order(start, end, filter, status, merchantID)

}
