package db

import "github.com/jmoiron/sqlx"

type GyftDB struct {
	*sqlx.DB
}

func NewDb(db *sqlx.DB) *GyftDB {
	return &GyftDB{db}
}
