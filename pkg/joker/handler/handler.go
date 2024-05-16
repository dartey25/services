package handler

import "github.com/mdoffice/md-services/pkg/eucustoms/service"

type JokerHandler struct {
	svc *service.EuCustomService
}

func New(s *service.EuCustomService) *JokerHandler {
	return &JokerHandler{svc: s}
}
