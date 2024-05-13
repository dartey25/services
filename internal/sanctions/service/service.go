package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

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

func (s *SanctionsService) Query(indexName string) (interface{}, error) {
	query := `{
  "query": {
    "fuzzy": {
      "name": {
        "value": "aviatsiinyi remontnyi zavod"
      }
    }
  }
}`
	res, err := s.es.Search(
		s.es.Search.WithContext(context.Background()),
		s.es.Search.WithIndex(indexName),
		s.es.Search.WithQuery(query),
		s.es.Search.WithTrackTotalHits(true),
		s.es.Search.WithHuman(),
		s.es.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%+v", res)

	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response interface{}
	err = json.Unmarshal(resBody, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
