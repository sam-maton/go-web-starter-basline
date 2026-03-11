package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	addr := flag.String("addr", ":4321", "HTTP network address")
	dsn := flag.String("dsn", "./sql/database.db", "SQLite db file location")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
