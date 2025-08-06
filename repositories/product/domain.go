package product

import (
	"time"

	dto "github.com/srv-cashpay/report/dto"
	"gorm.io/gorm"
)

type DomainReproductitory interface {
	GetReport(start, end time.Time, filter, status, merchantID string) ([]dto.DailyReportResponse, error)
}

type rproductReproductitory struct {
	DB *gorm.DB
}

func NewrproductReproductitory(DB *gorm.DB) DomainReproductitory {
	return &rproductReproductitory{
		DB: DB,
	}
}
