package main

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/yensho/get-your-fresh-tickets/api"
	"golang.org/x/exp/slog"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout))
	db := sqlx.MustConnect("postgres", "postgres://gyftusr:gyftPwE0@localhost:5432/postgres?sslmode=disable")

	api.NewApiServer(db, log).ListenAndServe()
	fmt.Println("Shutting down server...")
	os.Exit(1)
}
