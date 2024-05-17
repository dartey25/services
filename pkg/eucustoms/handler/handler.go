package handler

import (
	"github.com/mdoffice/md-services/pkg/eucustoms/service"
)

type EuCustomHandler struct {
	service *service.EuCustomService
}

func NewEuCustomHandler(s *service.EuCustomService) *EuCustomHandler {
	return &EuCustomHandler{service: s}
}
