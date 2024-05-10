package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mdoffice/md-services/model"
	"github.com/mdoffice/md-services/services"
	"github.com/mdoffice/md-services/views/templ/eucustom"
	"github.com/mdoffice/md-services/views/templ/search"
)

type EuCustomHandler struct {
	service *services.EuCustomService
}

func NewEuCustomHandler(s *services.EuCustomService) *EuCustomHandler {
	return &EuCustomHandler{service: s}
}

func (h *EuCustomHandler) HandleIndex(c echo.Context) error {
	return Render(c, http.StatusOK, eucustom.Index())
}

func (h *EuCustomHandler) HandleAeoTab(c echo.Context) error {
	queryTypes := c.QueryParams()["type"]

	types := make([]model.AeoType, 3)
	types[0] = model.AeoType{Code: "AEOC", Description: "Заява або дозвіл на отримання статусу Уповноваженого економічного оператора — Митні спрощення", Checked: false}
	types[1] = model.AeoType{Code: "AEOF", Description: "Заява або дозвіл на отримання статусу Уповноваженого економічного оператора — Митні спрощення/Безпека та захист", Checked: false}
	types[2] = model.AeoType{Code: "AEOS", Description: "Заява або дозвіл на отримання статусу Уповноваженого економічного оператора — Безпека та захист", Checked: false}

	for i := range types {
		for _, t := range queryTypes {
			if types[i].Code == t {
				types[i].Checked = true
			}
		}
	}

	return Render(c, http.StatusOK, eucustom.AeoTab(eucustom.AeoFormProps{Types: types}))
}

func (h *EuCustomHandler) HandleEoriTab(c echo.Context) error {
	return Render(c, http.StatusOK, eucustom.EoriTab())
}

func (h *EuCustomHandler) HandleAeoForm(c echo.Context) error {
	queryTypes := c.QueryParams()["type"]

	types := make([]model.AeoType, 3)
	types[0] = model.AeoType{Code: "AEOC", Description: "Заява або дозвіл на отримання статусу Уповноваженого економічного оператора — Митні спрощення", Checked: false}
	types[1] = model.AeoType{Code: "AEOF", Description: "Заява або дозвіл на отримання статусу Уповноваженого економічного оператора — Митні спрощення/Безпека та захист", Checked: false}
	types[2] = model.AeoType{Code: "AEOS", Description: "Заява або дозвіл на отримання статусу Уповноваженого економічного оператора — Безпека та захист", Checked: false}

	for i := range types {
		for _, t := range queryTypes {
			if types[i].Code == t {
				types[i].Checked = true
			}
		}
	}
	return Render(c, http.StatusOK, eucustom.AeoForm(eucustom.AeoFormProps{Types: types}))
}

func (h *EuCustomHandler) HandleEoriForm(c echo.Context) error {
	return Render(c, http.StatusOK, eucustom.EoriForm())
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
		return Render(c, http.StatusBadRequest, search.Error("at least one is required"))
	}

	results, err := h.service.GetAeoData(params)
	if err != nil {
		return Render(c, http.StatusInternalServerError, search.Error(err.Error()))
	}

	// if u := c.Request().Header.Get("HX-Current-URL"); u != "" {
	// 	if url, err := url.Parse(u); err == nil {
	// 		url.RawQuery = c.QueryParams().Encode()
	// 		c.Response().Header().Set("HX-Push-Url", url.String())
	// 	}
	// }

	if len(results.Data) == 0 {
		return Render(c, http.StatusNotFound, search.NotFound())
	}

	results.Page = params.Page
	results.Limit = params.Limit
	results.Query = strings.ToUpper(params.Holder)
	fmt.Printf("res: %v", results)
	if isHtmx {
		return Render(c, http.StatusOK, eucustom.AeoResults(results))
	}
	return c.JSON(http.StatusOK, results)
}

func (h *EuCustomHandler) HandleGetEoriData(c echo.Context) error {
	time.Sleep(time.Second * 4)
	eori := c.QueryParam("code")
	data, err := h.service.ValidateEori(eori)
	if err != nil {
		return Render(c, http.StatusOK, search.Error(err.Error()))
	}

	return Render(c, http.StatusOK, eucustom.EoriResults(model.NewEoriDisplayResult(data[0])))
}

func (h *EuCustomHandler) HandleJokerEoriData(c echo.Context) error {
	eori := c.QueryParam("code")
	data, err := h.service.ValidateEori(eori)
	if err != nil {
		return c.XML(http.StatusBadRequest, model.NewApiError(err.Error()))
	}

	return c.XML(http.StatusOK, model.NewApiResponse(data))
}
