package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Connect returns a pool of connection to the database, which is conncurrency
// safe. Uses the pgx interface.
func Connect(ctx context.Context, connstr string) (*pgxpool.Pool, error) {
	db, err := pgxpool.Connect(ctx, connstr)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to database: %w", err)
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error pinging database: %w", err)
	}

	return db, nil
}
