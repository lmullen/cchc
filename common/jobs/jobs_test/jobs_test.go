package items_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/lmullen/cchc/common/db"
	"github.com/lmullen/cchc/common/items"
	"github.com/lmullen/cchc/common/jobs"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJobsDB(t *testing.T) {
	t.Parallel()

	user := "gnomock"
	pass := "strong-passwords-are-the-best"
	dbname := "cchc_gnomock_test_jobs"

	p := postgres.Preset(
		postgres.WithUser(user, pass),
		postgres.WithDatabase(dbname),
	)

	container, err := gnomock.Start(p)
	assert.NoError(t, err)
	defer func() { assert.NoError(t, gnomock.Stop(container)) }()

	connstr := fmt.Sprintf("postgres://%s:%s@%s:%v/%s?sslmode=disable",
		user, pass, container.Host, container.DefaultPort(), dbname)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	db, _ := db.Connect(ctx, connstr)
	db.Ping(ctx)
	m, _ := migrate.New("file://../../../db/migrations", connstr)
	m.Up()

	var itemsRepo items.Repository
	itemsRepo = items.NewItemRepo(db)
	var jobsRepo jobs.Repository
	jobsRepo = jobs.NewJobsRepo(db)

	// Create an item from ID and URL
	item := &items.Item{
		ID: "http://www.loc.gov/item/mal1285100/",
		URL: sql.NullString{
			String: "https://www.loc.gov/item/mal1285100/",
			Valid:  true,
		},
	}

	// Save it to the repository. If we don't have an item in the repository,
	// then creating the jobs will fail the foreign key constraints.
	err = itemsRepo.Save(ctx, item)
	require.NoError(t, err)

	// Create a new job
	job := jobs.NewFullText(item.ID, "testing")

	// Save it to the database
	err = jobsRepo.SaveFullText(ctx, job)
	assert.NoError(t, err)

	job2, err := jobsRepo.GetFullText(ctx, job.ID)
	assert.NoError(t, err)
	assert.Equal(t, job, job2)

	// // Start and finish a job and save it to the database
	job.Start()
	job.Finish()
	err = jobsRepo.SaveFullText(ctx, job)
	assert.NoError(t, err)

}
