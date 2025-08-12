package product

import (
	"github.com/labstack/echo/v4"
	s "github.com/srv-cashpay/report/services/product"
)

type DomainHandler interface {
	BestSeller(c echo.Context) error
	SlowMoving(c echo.Context) error
}

type domainHandler struct {
	serviceRproduct s.RproductService
}

func NewRproductHandler(service s.RproductService) DomainHandler {
	return &domainHandler{
		serviceRproduct: service,
	}
}
