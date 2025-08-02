package pos

import (
	"github.com/labstack/echo/v4"
	s "github.com/srv-cashpay/report/services/pos"
)

type DomainHandler interface {
	Get(c echo.Context) error
}

type domainHandler struct {
	serviceRpos s.RposService
}

func NewRposHandler(service s.RposService) DomainHandler {
	return &domainHandler{
		serviceRpos: service,
	}
}
