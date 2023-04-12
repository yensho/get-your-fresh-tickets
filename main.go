package main

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/yensho/get-your-fresh-tickets/api"
	"golang.org/x/exp/slog"
	_ "modernc.org/sqlite"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout))
	db := sqlx.MustConnect("sqlite", "gyft.db")
	db.MustExec(`
	CREATE TABLE IF NOT EXISTS spaces (
		space_nm text NOT NULL PRIMARY KEY,
		space_section_nm text,
		space_section_seats blob
	);
	
	CREATE TABLE IF NOT EXISTS auth (
		user text NOT NULL PRIMARY KEY,
		token text NOT NULL
	)`)

	api.NewApiServer(db, log).ListenAndServe()
	fmt.Println("Shutting down server...")
	os.Exit(1)
}
