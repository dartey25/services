package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/mdoffice/md-services/internal/context"
	"github.com/mdoffice/md-services/pkg/core"
	"github.com/mdoffice/md-services/pkg/eucustoms/model"
)

func (h *JokerHandler) HandleEoriQuery(c echo.Context) error {
	eori := c.QueryParam("code")
	if eori == "" {
		eori = c.FormValue("code")
		if eori == "" {
			return c.XML(http.StatusOK, core.NewApiError("eori code is required, got none"))
		}
	}

	data, err := h.svc.eucustom.ValidateEori(eori)
	if err != nil {
		return c.XML(http.StatusOK, core.NewApiError(err.Error()))
	}

	return c.XML(http.StatusOK, core.NewApiResponse(data))
}

func parseAeoFormParams(c echo.Context, params *model.AeoQueryParams) {
	if limit, err := strconv.Atoi(c.FormValue("limit")); limit != 0 && err == nil {
		params.Limit = limit
	} else {
		params.Limit = 25
	}

	if page, err := strconv.Atoi(c.FormValue("page")); page != 0 && err == nil {
		params.Page = page
	} else {
		params.Page = 1
	}

	params.Holder = c.FormValue("holder")
	params.Country = c.FormValue("country")
	if p, err := c.FormParams(); err == nil {
		params.AuthTypes = p["type"]
	}
}

func (h *JokerHandler) HandleAeoQuery(c echo.Context) error {
	ctx := c.(*context.AppContext)
	var params model.AeoQueryParams
	parseAeoFormParams(c, &params)

	if params.Holder == "" && params.Country == "" && len(params.AuthTypes) == 0 {
		return ctx.XML(http.StatusBadRequest, "")
	}

	results, err := h.svc.eucustom.GetAeoData(params)
	if err != nil {
		return ctx.XML(http.StatusInternalServerError, "")
	}

	// if u := c.Request().Header.Get("HX-Current-URL"); u != "" {
	// 	if url, err := url.Parse(u); err == nil {
	// 		url.RawQuery = c.QueryParams().Encode()
	// 		c.Response().Header().Set("HX-Push-Url", url.String())
	// 	}
	// }

	if results == nil {
		return ctx.XML(http.StatusOK, "")
	}

	results.Page = params.Page
	results.Limit = params.Limit

	return ctx.XML(http.StatusOK, core.NewApiResponse(results))
}
