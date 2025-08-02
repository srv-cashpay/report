package pos

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/srv-cashpay/report/dto"
)

func (h *domainHandler) Get(c echo.Context) error {
	var filter dto.ReportFilterRequest
	if err := c.Bind(&filter); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid request"})
	}

	reports, err := h.serviceRpos.GetFilteredReport(filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Gagal mengambil data"})
	}

	return c.JSON(http.StatusOK, reports)
}
