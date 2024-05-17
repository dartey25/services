package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/mdoffice/md-services/internal/database"
	"github.com/mdoffice/md-services/pkg/sanctions/model"
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

func (s *SanctionsService) UploadInBatches(items []*model.SanctionsRow, indexName string, batchSize int) error {
	var (
		buf bytes.Buffer
		res *esapi.Response
		raw map[string]interface{}
		blk *database.BulkResponse

		numItems   int
		numErrors  int
		numIndexed int
		numBatches int
		currBatch  int
		count      int
	)

	count = len(items)

	if count%batchSize == 0 {
		numBatches = (count / batchSize)
	} else {
		numBatches = (count / batchSize) + 1
	}

	start := time.Now().UTC()
	for i, item := range items {
		numItems++

		currBatch = i / batchSize

		if i == count-1 {
			currBatch++
		}

		metaBytes := []byte(fmt.Sprintf(`{ "index" : { "_id" : "%d" } }%s`, item.SID, "\n"))
		data, err := json.Marshal(item)
		if err != nil {
			log.Fatalf("Cannot encode item with index %d: %s", i, err)
		}

		// Append newline to the data payload
		//
		data = append(data, "\n"...) // <-- Comment out to trigger failure for batch
		// fmt.Printf("%s", data) // <-- Uncomment to see the payload

		// // Uncomment next block to trigger indexing errors -->
		// if a.ID == 11 || a.ID == 101 {
		// 	data = []byte(`{"published" : "INCORRECT"}` + "\n")
		// }
		// // <--------------------------------------------------

		// Append payloads to the buffer (ignoring write errors)
		//
		buf.Grow(len(metaBytes) + len(data))
		buf.Write(metaBytes)
		buf.Write(data)

		// When a threshold is reached, execute the Bulk() request with body from buffer
		//
		if i > 0 && i%batchSize == 0 || i == count-1 {
			fmt.Printf("[%d/%d] ", currBatch, numBatches)

			res, err = s.es.Bulk(bytes.NewReader(buf.Bytes()), s.es.Bulk.WithIndex(indexName))
			if err != nil {
				log.Fatalf("Failure indexing batch %d: %s", currBatch, err)
				return fmt.Errorf("Failure indexing batch %d: %s", currBatch, err)
			}
			// If the whole request failed, print error and mark all documents as failed
			//
			if res.IsError() {
				numErrors += numItems
				if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
					log.Fatalf("Failure to to parse response body: %s", err)
					return fmt.Errorf("Failure to to parse response body: %s", err)
				} else {
					log.Printf("  Error: [%d] %s: %s",
						res.StatusCode,
						raw["error"].(map[string]interface{})["type"],
						raw["error"].(map[string]interface{})["reason"],
					)
					return fmt.Errorf("  Error: [%d] %s: %s",
						res.StatusCode,
						raw["error"].(map[string]interface{})["type"],
						raw["error"].(map[string]interface{})["reason"],
					)
				}
				// A successful response might still contain errors for particular documents...
				//
			} else {
				if err := json.NewDecoder(res.Body).Decode(&blk); err != nil {
					log.Fatalf("Failure to to parse response body: %s", err)
				} else {
					for _, d := range blk.Items {
						// ... so for any HTTP status above 201 ...
						//
						if d.Index.Status > 201 {
							// ... increment the error counter ...
							//
							numErrors++

							// ... and print the response status and error information ...
							log.Printf("  Error: [%d]: %s: %s: %s: %s",
								d.Index.Status,
								d.Index.Error.Type,
								d.Index.Error.Reason,
								d.Index.Error.Cause.Type,
								d.Index.Error.Cause.Reason,
							)
							return fmt.Errorf("  Error: [%d]: %s: %s: %s: %s",
								d.Index.Status,
								d.Index.Error.Type,
								d.Index.Error.Reason,
								d.Index.Error.Cause.Type,
								d.Index.Error.Cause.Reason,
							)
						} else {
							// ... otherwise increase the success counter.
							//
							numIndexed++
						}
					}
				}
			}

			// Close the response body, to prevent reaching the limit for goroutines or file handles
			//
			res.Body.Close()

			// Reset the buffer and items counter
			//
			buf.Reset()
			numItems = 0
		}
	}

	// Report the results: number of indexed docs, number of errors, duration, indexing rate
	//
	fmt.Print("\n")
	log.Println(strings.Repeat("▔", 65))

	dur := time.Since(start)

	if numErrors > 0 {
		log.Fatalf(
			"Indexed [%s] documents with [%s] errors in %s (%s docs/sec)",
			humanize.Comma(int64(numIndexed)),
			humanize.Comma(int64(numErrors)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(numIndexed))),
		)
		return fmt.Errorf(
			"Indexed [%s] documents with [%s] errors in %s (%s docs/sec)",
			humanize.Comma(int64(numIndexed)),
			humanize.Comma(int64(numErrors)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(numIndexed))),
		)
	} else {
		log.Printf(
			"Sucessfuly indexed [%s] documents in %s (%s docs/sec)",
			humanize.Comma(int64(numIndexed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(numIndexed))),
		)
	}

	return nil
}
