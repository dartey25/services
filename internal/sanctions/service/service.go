package service

import (
	"github.com/elastic/go-elasticsearch/v8"
)

type SanctionsService struct {
	es *elasticsearch.Client
}

func NewSanctionsService(e *elasticsearch.Client) *SanctionsService {
	return &SanctionsService{
		es: e,
	}
}
