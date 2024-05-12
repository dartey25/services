package db

import (
	"fmt"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/jmoiron/sqlx"
	"github.com/mdoffice/md-services/config"
	go_ora "github.com/sijms/go-ora/v2"
)

func NewOracleClient(cfg *config.DatabaseConfig) (db *sqlx.DB, err error) {
	connStr := go_ora.BuildUrl(cfg.Host, cfg.Port, cfg.Name, cfg.User, cfg.Password, nil)
	db, err = sqlx.Open("oracle", connStr)
	if err != nil {
		return
	}

	if err = db.Ping(); err != nil {
		return
	}

	return
}

func NewESClient(cfg *config.ElasticConfig) (*elasticsearch.Client, error) {
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

	return client, nil
}
