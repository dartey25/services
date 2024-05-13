package service

import (
	"github.com/jmoiron/sqlx"
)

type EuCustomService struct {
	db *sqlx.DB
}

func NewEuCustomService(db *sqlx.DB) *EuCustomService {
	return &EuCustomService{db: db}
}
