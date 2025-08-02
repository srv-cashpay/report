package pos

import (
	"time"

	dto "github.com/srv-cashpay/report/dto"
	"gorm.io/gorm"
)

type DomainRepository interface {
	GetReport(start, end time.Time) ([]dto.DailyReportResponse, error)
}

type rposRepository struct {
	DB *gorm.DB
}

func NewrposRepository(DB *gorm.DB) DomainRepository {
	return &rposRepository{
		DB: DB,
	}
}
