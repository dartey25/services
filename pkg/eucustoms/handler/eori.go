package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mdoffice/md-services/internal/context"
	"github.com/mdoffice/md-services/web/views/eucustom"
	"github.com/mdoffice/md-services/web/views/search"
)

func (h *EuCustomHandler) HandleEoriTab(c echo.Context) error {
	ctx := c.(context.AppContext)
	return ctx.Template(http.StatusOK, eucustom.EoriTab())
}

func (h *EuCustomHandler) HandleEoriForm(c echo.Context) error {
	ctx := c.(context.AppContext)
	return ctx.Template(http.StatusOK, eucustom.EoriForm())
}

func (h *EuCustomHandler) HandleGetEoriData(c echo.Context) error {
	ctx := c.(context.AppContext)
	eori := ctx.QueryParam("code")
	_, err := h.service.ValidateEori(eori)
	if err != nil {
		return ctx.Template(http.StatusInternalServerError, search.Error(err.Error()))
	}

	return ctx.Template(http.StatusInternalServerError, search.Error("Error message"))

	// return ctx.Template(http.StatusOK, eucustom.EoriResults(model.NewEoriDisplayResult(data[0])))
}
