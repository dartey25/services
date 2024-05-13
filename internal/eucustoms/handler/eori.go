package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mdoffice/md-services/internal/eucustoms/model"
	"github.com/mdoffice/md-services/pkg/core"
	"github.com/mdoffice/md-services/web/views/eucustom"
	"github.com/mdoffice/md-services/web/views/search"
)

func (h *EuCustomHandler) HandleEoriTab(c echo.Context) error {
	return core.Render(c, http.StatusOK, eucustom.EoriTab())
}

func (h *EuCustomHandler) HandleEoriForm(c echo.Context) error {
	return core.Render(c, http.StatusOK, eucustom.EoriForm())
}

func (h *EuCustomHandler) HandleGetEoriData(c echo.Context) error {
	time.Sleep(time.Second * 4)
	eori := c.QueryParam("code")
	data, err := h.service.ValidateEori(eori)
	if err != nil {
		return core.Render(c, http.StatusOK, search.Error(err.Error()))
	}

	return core.Render(c, http.StatusOK, eucustom.EoriResults(model.NewEoriDisplayResult(data[0])))
}

func (h *EuCustomHandler) HandleJokerEoriData(c echo.Context) error {
	eori := c.QueryParam("code")
	if eori == "" {
		return c.XML(http.StatusBadRequest, core.NewApiError("eori code is required, got none"))
	}
	data, err := h.service.ValidateEori(eori)
	if err != nil {
		return c.XML(http.StatusBadRequest, core.NewApiError(err.Error()))
	}

	return c.XML(http.StatusOK, core.NewApiResponse(data))
}
