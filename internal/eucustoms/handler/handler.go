package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mdoffice/md-services/internal/eucustoms/service"
	"github.com/mdoffice/md-services/pkg/core"
	"github.com/mdoffice/md-services/web/views/eucustom"
)

type EuCustomHandler struct {
	service *service.EuCustomService
}

func NewEuCustomHandler(s *service.EuCustomService) *EuCustomHandler {
	return &EuCustomHandler{service: s}
}

func (h *EuCustomHandler) HandleIndex(c echo.Context) error {
	return core.Render(c, http.StatusOK, eucustom.Index())
}
