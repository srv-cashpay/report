package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/srv-cashpay/middlewares/middlewares"
	"github.com/srv-cashpay/report/configs"

	h_pos "github.com/srv-cashpay/report/handlers/pos"
	r_pos "github.com/srv-cashpay/report/repositories/pos"
	s_pos "github.com/srv-cashpay/report/services/pos"

	h_product_best "github.com/srv-cashpay/report/handlers/product"
	r_product_best "github.com/srv-cashpay/report/repositories/product"
	s_product_best "github.com/srv-cashpay/report/services/product"
)

var (
	DB  = configs.InitDB()
	JWT = middlewares.NewJWTService()

	posR = r_pos.NewrposRepository(DB)
	posS = s_pos.NewRposService(posR, JWT)
	posH = h_pos.NewRposHandler(posS)

	product_bestR = r_product_best.NewrproductReproductitory(DB)
	product_bestS = s_product_best.NewRproductService(product_bestR, JWT)
	product_bestH = h_product_best.NewRproductHandler(product_bestS)
)

func New() *echo.Echo {

	e := echo.New()

	report := e.Group("/report", middlewares.AuthorizeJWT(JWT))
	{
		report.GET("/summary", posH.Summary)
		report.GET("/order", posH.Order)
		report.GET("/product_best", product_bestH.BestSeller)
		report.GET("/slow_moving", product_bestH.SlowMoving)

	}
	return e
}
