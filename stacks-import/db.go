package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// DBConnect connects to the database.
func DBConnect() (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	connstr, exists := os.LookupEnv("CCHC_DBSTR")
	if !exists {
		log.Fatal("Database connection string not set as an environment variable")
	}

	db, err := pgxpool.Connect(ctx, connstr)
	if err != nil {
		return nil, err
	}
	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return db, err
}
