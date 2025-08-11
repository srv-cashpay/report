package pos

import (
	"time"

	"github.com/srv-cashpay/report/dto"
)

func (s *rposService) Summary(start, end time.Time, filter, status, merchantID string) ([]dto.DailyReportResponse, error) {
	return s.Repo.Summary(start, end, filter, status, merchantID)

}
