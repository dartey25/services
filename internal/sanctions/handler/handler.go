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
	return nil
}

func (h *SanctionsHandler) HandleQueryLegal(c echo.Context) error {
	return nil
}
