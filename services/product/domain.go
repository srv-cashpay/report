package product

import (
	"time"

	m "github.com/srv-cashpay/middlewares/middlewares"

	"github.com/srv-cashpay/report/dto"
	r "github.com/srv-cashpay/report/repositories/product"
)

type RproductService interface {
	BestSeller(start, end time.Time, filter, status, merchantID string) ([]dto.BestSellerResponse, error)
}

type rproductService struct {
	Repo r.DomainReproductitory
	jwt  m.JWTService
}

func NewRproductService(Repo r.DomainReproductitory, jwtS m.JWTService) RproductService {
	return &rproductService{
		Repo: Repo,
		jwt:  jwtS,
	}
}
