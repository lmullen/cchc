package db

import (
	"context"
	"embed"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" //Driver for migrations
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrations embed.FS

// MigrateUp brings the database to the current schema
func MigrateUp(ctx context.Context, connstr string) error {
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

// MigrateDown resets the database back to a blank state. USE WITH
// CAUTION: THis will result in data loss.
func MigrateDown(ctx context.Context, connstr string) error {
	d, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}
	m, err := migrate.NewWithSourceInstance("iofs", d, connstr)
	if err != nil {
		return err
	}
	err = m.Down()
	if err != nil {
		return err
	}

	return nil
}
