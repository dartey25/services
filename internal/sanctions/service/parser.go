package service

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/mdoffice/md-services/internal/sanctions/model"
	"github.com/xuri/excelize/v2"
)

func fallible(a, b int) int {
	if a > -1 {
		return a
	}
	return b
}

func parseId(s string) *model.Identificator {
	var id model.Identificator
	idxT := strings.Index(s, ":")
	if idxT != -1 {
		str := s[:idxT]
		id.Type = &str
	}

	idx := strings.Index(s, "[")
	str := s[fallible(idxT+1, 0):fallible(idx, len(s))]
	id.Value = &str

	return &id
}

func (s *SanctionsService) ParseLegal(filePath string) ([]*model.SanctionsRow, error) {
	timeBegin := time.Now()
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	}()

	var sanctionRows []*model.SanctionsRow
	for _, sheetName := range f.GetSheetList() {
		rows, err := f.GetRows(sheetName)
		if err != nil {
			return nil, err
		}
		sanctionRows = append(sanctionRows, parseRows(rows)...)
	}

	timeEnd := time.Now()
	log.Printf("Execution time %s, items parced: %d", timeEnd.Sub(timeBegin), len(sanctionRows))
	return sanctionRows, nil
}

func parseMap(s string) (map[string]string, bool) {
	m := make(map[string]string)
	for _, s := range strings.Split(s, ";") {
		kh := strings.Split(strings.TrimSpace(s), "]")
		if len(kh) < 2 {
			return nil, false
		}
		key := strings.TrimPrefix(kh[0], "[")
		value := strings.TrimSpace(kh[1])
		m[key] = value
	}

	return m, true
}

func parseRows(r [][]string) []*model.SanctionsRow {
	rows := make([]*model.SanctionsRow, 0)
	for i, row := range r {
		if i == 0 {
			continue
		}

		var sRow model.SanctionsRow
		for j, v := range row {
			switch j {
			case 0:
				sid, err := strconv.Atoi(v)
				if err != nil {
					// TODO: log err
					break
				}
				sRow.SID = sid
			case 1:
				sRow.Company.Name = v
			case 2:
				sRow.Company.TranslitName = v
			case 3:
				sRow.Company.Aliases = strings.Split(v, "; ")
			case 4:
				sRow.Company.Status = v
			case 5:
				sRow.Company.Country = v
			case 6:
				sRow.Company.RegID = parseId(v)
			case 7:
				sRow.Company.TaxID = parseId(v)
			case 8:
				sRow.Company.AdditionalInfo = strings.Split(v, "; ")
			case 9:
				if m, ok := parseMap(v); ok {
					sRow.Sanctions.Active = m
				}
			case 10:
				sRow.Sanctions.Term = v
			case 11:
				time, err := time.Parse("2006-01-02", v)
				if err != nil {
					sRow.Sanctions.EndDate = nil
				}
				sRow.Sanctions.EndDate = &time
			case 12:
				sRow.LastAction.Type = v
			case 13:
				sRow.LastAction.Document.Num = v
			case 14:
				time, err := time.Parse("2006-01-02", v)
				if err != nil {
					sRow.LastAction.Document.Date = nil
				}
				sRow.LastAction.Document.Date = &time

			case 15:
				sRow.LastAction.Document.Appendix = v
			case 16:
				sRow.LastAction.Document.AppendixPosition = v
			}
		}
		rows = append(rows, &sRow)
	}

	return rows
}

func (s *SanctionsService) ReCreateIndex(indexName string) error {
	if _, err := s.es.Indices.Delete([]string{indexName}); err != nil {
		// log.Fatalf("Cannot delete index: %s", err)
		return err
	}
	res, err := s.es.Indices.Create(indexName)
	if err != nil {
		// log.Fatalf("Cannot create index: %s", err)
		return err
	}
	if res.IsError() {
		// log.Fatalf("Cannot create index: %s", res)
		return err
	}

	return nil
}

func (s *SanctionsService) UploadInBatches(items []*model.SanctionsRow, indexName string) error {
	return nil
}
