package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/mdoffice/md-services/config"
	go_ora "github.com/sijms/go-ora/v2"
)

func NewOracleClient(cfg *config.DatabaseConfig) (db *sqlx.DB, err error) {
	connStr := go_ora.BuildUrl(cfg.Host, cfg.Port, cfg.Name, cfg.User, cfg.Password, nil)
	fmt.Printf("Connection string: %s\n", connStr)
	db, err = sqlx.Open("oracle", connStr)
	if err != nil {
		return
	}

	if err = db.Ping(); err != nil {
		return
	}

	return
}
