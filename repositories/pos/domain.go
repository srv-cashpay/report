package pos

import (
	"time"

	dto "github.com/srv-cashpay/report/dto"
	"gorm.io/gorm"
)

type DomainRepository interface {
	Summary(start, end time.Time, filter, status, merchantID string) ([]dto.DailyReportResponse, error)
	Order(start, end time.Time, filter, status, merchantID string) ([]dto.DailyReportResponse, error)
}

type rposRepository struct {
	DB *gorm.DB
}

func NewrposRepository(DB *gorm.DB) DomainRepository {
	return &rposRepository{
		DB: DB,
	}
}
