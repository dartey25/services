package handler

import (
	"github.com/mdoffice/md-services/pkg/eucustoms/service"
)

type Services struct {
	eucustom *service.EuCustomService
}

type JokerHandler struct {
	svc Services
}

func New(s *service.EuCustomService) *JokerHandler {
	services := Services{
		eucustom: s,
	}
	return &JokerHandler{svc: services}
}
