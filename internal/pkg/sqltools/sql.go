package sqltools

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"db-performance-project/internal/pkg"

	// justifying it
	_ "github.com/jackc/pgx/stdlib"
)

func InsertBatch(ctx context.Context, db *sql.DB, query string, values []interface{}) (*sql.Rows, error) {
	rows, err := db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, fmt.Errorf("InsertBatch: [%w] when inserting row into [%s] table \n [%+v]", err, query, values)
	}

	return rows, nil
}

func NewPostgresURL() string {
	url := "user=" + os.Getenv("POSTGRES_USER") +
		" dbname=" + os.Getenv("POSTGRES_DB") +
		" password=" + os.Getenv("POSTGRES_PASSWORD") +
		" host=" + os.Getenv("POSTGRES_HOST") +
		" port=" + os.Getenv("POSTGRES_PORT") +
		" sslmode=" + os.Getenv("POSTGRES_SSLMODE")

	return url
}

type Database struct {
	Connection *sql.DB
}

func NewPostgresRepository(config *pkg.DatabaseParams) *Database {
	connection, err := sql.Open("pgx", NewPostgresURL())
	if err != nil {
		log.Fatalln("Can't parse config", err)
	}

	connection.SetMaxOpenConns(config.MaxOpenCons)

	err = connection.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	return &Database{Connection: connection}
}
