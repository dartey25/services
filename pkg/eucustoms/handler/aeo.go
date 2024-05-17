package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/mdoffice/md-services/pkg/core"
	"github.com/mdoffice/md-services/pkg/eucustoms/model"
	"github.com/mdoffice/md-services/web/views/eucustom"
	"github.com/mdoffice/md-services/web/views/search"
)

func getAuthTypes(queryTypes []string) []*model.AeoType {
	types := make([]*model.AeoType, 3)
	types[0] = &model.AeoType{Code: "AEOC", Description: "Заява або дозвіл на отримання статусу Уповноваженого економічного оператора — Митні спрощення", Checked: true}
	types[1] = &model.AeoType{Code: "AEOF", Description: "Заява або дозвіл на отримання статусу Уповноваженого економічного оператора — Митні спрощення/Безпека та захист", Checked: true}
	types[2] = &model.AeoType{Code: "AEOS", Description: "Заява або дозвіл на отримання статусу Уповноваженого економічного оператора — Безпека та захист", Checked: true}

	for i := range types {
		for _, t := range queryTypes {
			if types[i].Code != t {
				types[i].Checked = false
			}
		}
	}
	return types
}

func (h *EuCustomHandler) HandleIndex(c echo.Context) error {
	queryTypes := c.QueryParams()["type"]
	types := getAuthTypes(queryTypes)
	countries, err := h.service.GetCountries()
	if err != nil {
		return err
	}

	return core.Render(c, http.StatusOK, eucustom.Index(eucustom.AeoFormProps{Types: types, CountryList: countries}))
}

func (h *EuCustomHandler) HandleAeoTab(c echo.Context) error {
	queryTypes := c.QueryParams()["type"]
	types := getAuthTypes(queryTypes)
	countries, err := h.service.GetCountries()
	if err != nil {
		return err
	}
	return core.Render(c, http.StatusOK, eucustom.AeoTab(eucustom.AeoFormProps{Types: types, CountryList: countries}))
}

func (h *EuCustomHandler) HandleAeoForm(c echo.Context) error {
	queryTypes := c.QueryParams()["type"]
	types := getAuthTypes(queryTypes)
	countries, err := h.service.GetCountries()
	if err != nil {
		return err
	}

	return core.Render(c, http.StatusOK, eucustom.AeoForm(eucustom.AeoFormProps{Types: types, CountryList: countries}))
}

func parseAeoQueryParams(c echo.Context, params *model.AeoQueryParams) {
	if limit, err := strconv.Atoi(c.QueryParam("limit")); limit != 0 && err == nil {
		params.Limit = limit
	} else {
		params.Limit = 25
	}

	if page, err := strconv.Atoi(c.QueryParam("page")); page != 0 && err == nil {
		params.Page = page
	} else {
		params.Page = 1
	}

	params.Holder = c.QueryParam("holder")
	params.Country = c.QueryParam("country")
	params.AuthTypes = c.QueryParams()["type"]
}

func (h *EuCustomHandler) HandleGetAeoData(c echo.Context) error {
	isHtmx := c.Request().Header.Get("HX-Request") == "true"
	var params model.AeoQueryParams
	parseAeoQueryParams(c, &params)

	if params.Holder == "" && params.Country == "" && len(params.AuthTypes) == 0 {
		return core.Render(c, http.StatusBadRequest, search.Error("at least one is required"))
	}

	results, err := h.service.GetAeoData(params)
	if err != nil {
		return core.Render(c, http.StatusInternalServerError, search.Error(err.Error()))
	}

	// if u := c.Request().Header.Get("HX-Current-URL"); u != "" {
	// 	if url, err := url.Parse(u); err == nil {
	// 		url.RawQuery = c.QueryParams().Encode()
	// 		c.Response().Header().Set("HX-Push-Url", url.String())
	// 	}
	// }

	if len(results.Data) == 0 {
		return core.Render(c, http.StatusNotFound, search.NotFound())
	}

	results.Page = params.Page
	results.Limit = params.Limit
	if isHtmx {
		return core.Render(c, http.StatusOK, eucustom.AeoResults(results))
	}
	return c.JSON(http.StatusOK, results)
}
