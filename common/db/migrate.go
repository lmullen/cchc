package db

import (
	"context"
	"embed"

	"github.com/golang-migrate/migrate/v4" // Driver for migrations
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrations embed.FS

// Migrate brings the database to the current schema
func Migrate(ctx context.Context, connstr string) error {
	d, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}
	m, err := migrate.NewWithSourceInstance("iofs", d, connstr)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil {
		return err
	}

	return nil
}
