package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mdoffice/md-services/internal/sanctions/service"
)

type SanctionsHandler struct {
	svc *service.SanctionsService
}

func NewSanctionsHandler(s *service.SanctionsService) *SanctionsHandler {
	return &SanctionsHandler{
		svc: s,
	}
}

const LEGAL_SANCTIONS_IDX = "legal_sanctions"

func (h *SanctionsHandler) HandleParseLegal(c echo.Context) error {
	sanctions, err := h.svc.ParseLegal("/home/dartey/projects/mdoffice/sanctions-parser/legal_sanctions.xlsx")
	if err != nil {
		c.Logger().Errorf("Error parsing sanctions file: %s", err)
		return err
	}

	c.Logger().Info("Parsed file")

	err = h.svc.ReCreateIndex(LEGAL_SANCTIONS_IDX)
	if err != nil {
		return err
	}
	c.Logger().Info("Recreated index file")

	err = h.svc.UploadInBatches(sanctions, LEGAL_SANCTIONS_IDX, 1000)
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, "Successful upload")
}

func (h *SanctionsHandler) HandleQueryLegal(c echo.Context) error {
	res, err := h.svc.Query(LEGAL_SANCTIONS_IDX)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}
