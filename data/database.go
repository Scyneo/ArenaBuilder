package data

import (
  "context"
  "fmt"
  "github.com/jackc/pgx/v5/pgxpool"
  "log"
  "os"
)

func NewDb() (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context.Background(), "postgresql://postgres:admin@localhost:5432/lol")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
    return nil, err
	}
  log.Println("Connected to database")

  return dbpool, nil
}
