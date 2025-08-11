package product

import (
	"time"

	"github.com/srv-cashpay/report/dto"
)

func (s *rproductService) BestSeller(start, end time.Time, filter, status, merchantID string) ([]dto.BestSellerResponse, error) {
	return s.Repo.BestSeller(start, end, filter, status, merchantID)

}
