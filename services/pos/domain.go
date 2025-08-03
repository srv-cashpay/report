package pos

import (
	"time"

	m "github.com/srv-cashpay/middlewares/middlewares"

	"github.com/srv-cashpay/report/dto"
	r "github.com/srv-cashpay/report/repositories/pos"
)

type RposService interface {
	GetFilteredReport(start, end time.Time, filter, status, merchantID string) ([]dto.DailyReportResponse, error)
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
