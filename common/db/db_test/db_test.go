package db_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lmullen/cchc/common/db"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/postgres"
	"github.com/stretchr/testify/require"
)

func TestDBConnection(t *testing.T) {
	t.Parallel()

	user := "gnomock"
	pass := "strong-passwords-are-the-best"
	dbname := "cchc_gnomock_test"

	p := postgres.Preset(
		postgres.WithUser(user, pass),
		postgres.WithDatabase(dbname),
	)

	container, err := gnomock.Start(p)
	require.NoError(t, err)
	defer func() { require.NoError(t, gnomock.Stop(container)) }()

	connstr := fmt.Sprintf("postgres://%s:%s@%s:%v/%s?sslmode=disable",
		user, pass, container.Host, container.DefaultPort(), dbname)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := db.Connect(ctx, connstr)
	require.NoError(t, err)

	err = db.Ping(ctx)
	require.NoError(t, err)

	require.IsType(t, &pgxpool.Pool{}, db)

	m, err := migrate.New("file://../../../db/migrations", connstr)
	require.NoError(t, err)

	err = m.Up()
	require.NoError(t, err)

	err = m.Down()
	require.NoError(t, err)

}
