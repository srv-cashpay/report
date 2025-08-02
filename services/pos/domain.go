package pos

import (
	m "github.com/srv-cashpay/middlewares/middlewares"

	"github.com/srv-cashpay/report/dto"
	r "github.com/srv-cashpay/report/repositories/pos"
)

type RposService interface {
	GetFilteredReport(filter dto.ReportFilterRequest) ([]dto.DailyReportResponse, error)
}

type rposService struct {
	Repo r.DomainRepository
	jwt  m.JWTService
}

func NewRposService(Repo r.DomainRepository, jwtS m.JWTService) RposService {
	return &rposService{
		Repo: Repo,
		jwt:  jwtS,
	}
}
