package db

import (
	"time"

	"github.com/google/uuid"
)

type Space struct {
	Name    string `db:"space_nm"`
	Section string `db:"space_section_nm"`
	Seats   []byte `db:"space_section_seats"`
}

type Customer struct {
	ID         uuid.UUID              `db:"id"`
	Name       string                 `db:"cust_name"`
	Age        int                    `db:"age"`
	Gender     string                 `db:"gender"`
	PhoneNum   string                 `db:"phone_num"`
	Address    map[string]interface{} `db:"addr"`
	Email      string                 `db:"email"`
	InsertTime time.Time              `db:"isrt_ts"`
	InsertUser string                 `db:"isrt_usr"`
	UpdateTime time.Time              `db:"updt_ts"`
	UpdateUser string                 `db:"updt_usr"`
}

type Show struct {
	ID          uuid.UUID              `db:"id"`
	Name        string                 `db:"show_name"`
	Location    string                 `db:"loc"`
	ShowTime    time.Time              `db:"show_time"`
	Seats       []string               `db:"seats"`
	Prices      map[string]interface{} `db:"prices"`
	AgeRestrict bool                   `db:"age_restrict"`
	InsertTime  time.Time              `db:"isrt_ts"`
	InsertUser  string                 `db:"isrt_usr"`
	UpdateTime  time.Time              `db:"updt_ts"`
	UpdateUser  string                 `db:"updt_usr"`
}
