package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/srv-cashpay/middlewares/middlewares"
	"github.com/srv-cashpay/report/configs"

	h_pos "github.com/srv-cashpay/report/handlers/pos"
	r_pos "github.com/srv-cashpay/report/repositories/pos"
	s_pos "github.com/srv-cashpay/report/services/pos"
)

var (
	DB  = configs.InitDB()
	JWT = middlewares.NewJWTService()

	posR = r_pos.NewrposRepository(DB)
	posS = s_pos.NewRposService(posR, JWT)
	posH = h_pos.NewRposHandler(posS)
)

func New() *echo.Echo {

	e := echo.New()

	report := e.Group("api/report", middlewares.AuthorizeJWT(JWT))
	{
		report.GET("/pos", posH.Get)
	}
	return e
}
