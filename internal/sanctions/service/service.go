package service

import (
	"github.com/mdoffice/md-services/internal/db"
)

type SanctionsService struct {
	es *db.ElasticClient
}

func NewSanctionsService(e *db.ElasticClient) *SanctionsService {
	return &SanctionsService{
		es: e,
	}
}
