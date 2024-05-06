package handlers

import (
	"net/http"
	"net/url"
	"strconv"

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
	types[0] = model.AeoType{Code: "AEOC", Description: "Application or authorisation for the status of Authorised Economic Operator — Customs simplifications", Checked: false}
	types[1] = model.AeoType{Code: "AEOF", Description: "Application or authorisation for the status of Authorised Economic Operator — Customs simplifications/Security and safety", Checked: false}
	types[2] = model.AeoType{Code: "AEOS", Description: "Application or authorisation for the status of Authorised Economic Operator — Security and safety", Checked: false}

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
	types[0] = model.AeoType{Code: "AEOC", Description: "Application or authorisation for the status of Authorised Economic Operator — Customs simplifications", Checked: false}
	types[1] = model.AeoType{Code: "AEOF", Description: "Application or authorisation for the status of Authorised Economic Operator — Customs simplifications/Security and safety", Checked: false}
	types[2] = model.AeoType{Code: "AEOS", Description: "Application or authorisation for the status of Authorised Economic Operator — Security and safety", Checked: false}

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

func (h *EuCustomHandler) HandleGetAeoData(c echo.Context) error {
	queryParams := c.QueryParams()
	holder := c.QueryParam("holder")
	country := c.QueryParam("country")
	types := c.QueryParams()["type"]
	const limit = 25
	var page int
	var err error
	isHtmx := c.Request().Header.Get("HX-Request") == "true"

	if pageStr := c.QueryParam("page"); pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			return err
		}
	}

	if holder == "" && country == "" && len(types) == 0 {
		return Render(c, http.StatusBadRequest, search.Error("at least one is required"))
	}

	results, err := h.service.GetAeoData(holder, country, types, page, limit)
	if err != nil {
		return Render(c, http.StatusInternalServerError, search.Error(err.Error()))
	}

	if u := c.Request().Header.Get("HX-Current-URL"); u != "" {
		if url, err := url.Parse(u); err == nil {
			url.RawQuery = queryParams.Encode()
			c.Response().Header().Set("HX-Push-Url", url.String())
		}
	}

	if len(results.Data) == 0 {
		return Render(c, http.StatusNotFound, search.NotFound())
	}

	results.NextPage = page + 1
	if isHtmx {
		return Render(c, http.StatusOK, search.AeoResults(results))
	}
	return c.JSON(http.StatusOK, results)
}

func (h *EuCustomHandler) HandleGetEoriData(c echo.Context) error {
	eori := c.QueryParam("code")
	data, err := h.service.ValidateEori(eori)
	if err != nil {
		return Render(c, http.StatusOK, search.Error(err.Error()))
	}

	return Render(c, http.StatusOK, search.EoriResults(model.NewEoriDisplayResult(data[0])))
}

func (h *EuCustomHandler) HandleJokerEoriData(c echo.Context) error {
	eori := c.QueryParam("code")
	data, err := h.service.ValidateEori(eori)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, model.NewApiResponse(data))
}
