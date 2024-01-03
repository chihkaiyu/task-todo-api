package postgres

import (
	"runtime"
	"time"

	"github.com/jmoiron/sqlx"

	// init driver
	_ "github.com/lib/pq"
)

func New(datasource string) (*sqlx.DB, error) {
	dbx, err := sqlx.Connect("postgres", datasource)
	if err != nil {
		return nil, err
	}

	dbx.DB.SetConnMaxLifetime(10 * time.Second)
	dbx.DB.SetMaxOpenConns(runtime.GOMAXPROCS(0) * 2)

	return dbx, nil
}
