package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {
	// localPostgresURL := "postgresql://root:admin@localhost:5433/whisper-warp-db?sslmode=disable"

	POSTGRES_URL := os.Getenv("POSTGRES_URL")
	fmt.Println("POSTGRES_URL: ", POSTGRES_URL)

	db, err := sql.Open("postgres", POSTGRES_URL)
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}
