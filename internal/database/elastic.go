package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/mdoffice/md-services/config"
)

type ElasticClient struct {
	*elasticsearch.Client
}

func NewElasticClient(cfg *config.ElasticConfig) (*ElasticClient, error) {
	fmt.Println(cfg.CertPath)
	certBytes, err := os.ReadFile(cfg.CertPath)
	if err != nil {
		return nil, fmt.Errorf("error reading http certificate: %s", err)
	}

	url := fmt.Sprintf("https://%s:%d/", cfg.Host, cfg.Port)

	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			url,
		},
		APIKey: cfg.ApiKey,
		CACert: certBytes,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating the client: %s", err)
	}

	return &ElasticClient{client}, nil
}

type ElasticResponse struct {
	*esapi.Response
}

type BulkResponse struct {
	Errors bool `json:"errors"`
	Items  []struct {
		Index struct {
			ID     string `json:"_id"`
			Result string `json:"result"`
			Status int    `json:"status"`
			Error  struct {
				Type   string `json:"type"`
				Reason string `json:"reason"`
				Cause  struct {
					Type   string `json:"type"`
					Reason string `json:"reason"`
				} `json:"caused_by"`
			} `json:"error"`
		} `json:"index"`
	} `json:"items"`
}

type BulkUploadParams struct {
	Items []struct {
		Meta string
		Data interface{}
	}
	IndexName string
	BatchSize int
}

func (c *ElasticClient) BulkUpload(params *BulkUploadParams) error {
	var (
		buf bytes.Buffer
		res *esapi.Response
		raw map[string]interface{}
		blk *BulkResponse

		numItems   int
		numErrors  int
		numIndexed int
		numBatches int
		currBatch  int
		count      int
	)

	count = len(params.Items)

	if count%params.BatchSize == 0 {
		numBatches = (count / params.BatchSize)
	} else {
		numBatches = (count / params.BatchSize) + 1
	}

	start := time.Now().UTC()

	// Loop over the collection
	//
	for i, item := range params.Items {
		numItems++

		currBatch = i / params.BatchSize

		if i == count-1 {
			currBatch++
		}

		metaBytes := []byte(item.Meta)
		data, err := json.Marshal(item.Data)
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
		if i > 0 && i%params.BatchSize == 0 || i == count-1 {
			fmt.Printf("[%d/%d] ", currBatch, numBatches)

			res, err = c.Bulk(bytes.NewReader(buf.Bytes()), c.Bulk.WithIndex(params.IndexName))
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
