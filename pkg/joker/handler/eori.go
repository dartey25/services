package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mdoffice/md-services/pkg/core"
)

func (h *JokerHandler) HandleEoriQuery(c echo.Context) error {
	eori := c.QueryParam("code")
	if eori == "" {
		eori = c.FormValue("code")
		if eori == "" {
			return c.XML(http.StatusOK, core.NewApiError("eori code is required, got none"))
		}
	}

	data, err := h.svc.ValidateEori(eori)
	if err != nil {
		return c.XML(http.StatusOK, core.NewApiError(err.Error()))
	}

	return c.XML(http.StatusOK, core.NewApiResponse(data))
}
