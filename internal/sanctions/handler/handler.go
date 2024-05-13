package handler

import (
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

func (h *SanctionsHandler) HandleParseLegal(c echo.Context) error {
	sanctions, err := h.svc.ParseLegal("path_to_file")
	if err != nil {
		c.Logger().Errorf("Error parsing sanctions file: %s", err)
		return err
	}

	err = h.svc.ReCreateIndex("index_name")
	if err != nil {
		return err
	}

	err = h.svc.UploadInBatches(sanctions, "index_name")
	if err != nil {
		return err
	}

	return nil
}

func (h *SanctionsHandler) HandleQueryLegal(c echo.Context) error {
	return nil
}
