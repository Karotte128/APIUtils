package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ConnPool *pgxpool.Pool

func CreateConnection(pgxconf string) {
	poolConfig, err := pgxpool.ParseConfig(pgxconf)
	ConnPool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	err = ConnPool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Ping to database failed: %v\n", err)
	}
}
