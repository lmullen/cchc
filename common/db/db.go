package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/cenkalti/backoff/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Connect returns a pool of connection to the database, which is conncurrency
// safe. Uses the pgx interface.
func Connect(ctx context.Context, connstr string, application string) (*pgxpool.Pool, error) {

	connstr, err := AddApplication(connstr, application)
	if err != nil {
		return nil, fmt.Errorf("Error adding application name to connection string: %w", err)
	}

	var db *pgxpool.Pool

	connectWithRetry := func() error {
		select {
		case <-ctx.Done():
			return backoff.Permanent(errors.New("Cancelled attempt to connect to database"))
		default:
			conn, err := pgxpool.Connect(ctx, connstr)
			if err != nil {
				return fmt.Errorf("Error connecting to database: %w", err)
			}

			err = conn.Ping(ctx)
			if err != nil {
				return fmt.Errorf("Error pinging database: %w", err)
			}
			db = conn
		}
		return nil
	}

	policy := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 10)
	err = backoff.Retry(connectWithRetry, policy)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to database: %w", err)
	}

	return db, nil
}
