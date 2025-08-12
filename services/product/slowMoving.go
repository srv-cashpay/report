package product

import (
	"time"

	"github.com/srv-cashpay/report/dto"
)

func (s *rproductService) SlowMoving(start, end time.Time, filter, status, merchantID string) ([]dto.BestSellerResponse, error) {
	return s.Repo.SlowMoving(start, end, filter, status, merchantID)

}
