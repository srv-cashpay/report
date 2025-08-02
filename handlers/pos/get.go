package pos

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/srv-cashpay/report/dto"
)

func (h *domainHandler) Get(c echo.Context) error {
	var req dto.ReportFilterRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request params")
	}

	start, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid start_date format, must be YYYY-MM-DD")
	}

	end, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid end_date format, must be YYYY-MM-DD")
	}

	if req.Filter == "" {
		req.Filter = "daily"
	}

	data, err := h.serviceRpos.GetFilteredReport(start, end, req.Filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, data)
}
