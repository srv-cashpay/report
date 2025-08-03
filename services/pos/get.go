package pos

import (
	"time"

	"github.com/srv-cashpay/report/dto"
)

// func (s *rposService) GetFilteredReport(req dto.ReportFilterRequest) ([]dto.DailyReportResponse, error) {
// 	layout := "2006-01-02"
// 	start, _ := time.Parse(layout, req.StartDate)
// 	end, _ := time.Parse(layout, req.EndDate)

// 	return s.Repo.GetReport(start, end)
// }

func (s *rposService) GetFilteredReport(start, end time.Time, filter, status, merchantID string) ([]dto.DailyReportResponse, error) {
	return s.Repo.GetReport(start, end, filter, status, merchantID)

}
