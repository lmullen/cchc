package items_test

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/mpvl/unique"

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

	db, _ := db.Connect(ctx, connstr, "jobs-test")
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

func TestJobsStatus(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	job := jobs.NewFullText("test_item", "test_queue")
	assert.Equal("ready", job.Status)

	job.Start()
	assert.Equal("running", job.Status)

	job.Skip()
	assert.Equal("skipped", job.Status)

	job.Fail()
	assert.Equal("failed", job.Status)

	job.Finish()
	assert.Equal("finished", job.Status)
}

// Test that we can enqueue jobs without errors
func TestEnqueingJobs(t *testing.T) {
	t.Parallel()

	user := "gnomock"
	pass := "strong-passwords-are-the-best"
	dbname := "cchc_gnomock_test_job_queue"

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

	db, _ := db.Connect(ctx, connstr, "jobs-test")
	db.Ping(ctx)
	m, _ := migrate.New("file://../../../db/migrations", connstr)
	m.Up()

	var itemsRepo items.Repository
	itemsRepo = items.NewItemRepo(db)
	var jobsRepo jobs.Repository
	jobsRepo = jobs.NewJobsRepo(db)

	totaltests := 500
	halftests := 250

	for i := 0; i < totaltests; i++ {
		item := &items.Item{
			ID:  fmt.Sprintf("%04d", i),
			API: sql.NullString{String: "{\"test\":\"test\"}", Valid: true},
		}
		err := itemsRepo.Save(ctx, item)
		require.NoError(t, err)
	}

	wg := &sync.WaitGroup{}
	var items1 []string
	var items2 []string

	testQueue := func(items *[]string) {
		defer wg.Done()
		for i := 0; i < halftests; i++ {
			job, err := jobsRepo.CreateJobForUnqueued(ctx, "testing")
			require.NoError(t, err)
			*items = append(*items, job.ItemID)
		}
	}

	wg.Add(1)
	go testQueue(&items1)
	wg.Add(1)
	go testQueue(&items2)
	wg.Wait()

	allItems := append(items1, items2...)

	assert.Equal(t, totaltests, len(allItems))

	sort.Strings(allItems)
	assert.True(t, unique.StringsAreUnique(allItems))

}

func TestNoJobsNeedEnqueuing(t *testing.T) {
	t.Parallel()

	user := "gnomock"
	pass := "strong-passwords-are-the-best"
	dbname := "cchc_gnomock_test_job_queue_no_items"

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

	db, _ := db.Connect(ctx, connstr, "jobs-test")
	db.Ping(ctx)
	m, _ := migrate.New("file://../../../db/migrations", connstr)
	m.Up()

	var jobsRepo jobs.Repository
	jobsRepo = jobs.NewJobsRepo(db)

	// There are no items in the database, so we should expect not to make any jobs.
	job, err := jobsRepo.CreateJobForUnqueued(ctx, "testing")

	assert.Nil(t, job)
	assert.ErrorIs(t, err, jobs.ErrAllQueued)

}

func TestGettingJobsFromQueue(t *testing.T) {
	t.Parallel()

	user := "gnomock"
	pass := "strong-passwords-are-the-best"
	dbname := "cchc_gnomock_test_job_queue_retrieval"

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

	db, _ := db.Connect(ctx, connstr, "jobs-test")
	db.Ping(ctx)
	m, _ := migrate.New("file://../../../db/migrations", connstr)
	m.Up()

	var itemsRepo items.Repository
	itemsRepo = items.NewItemRepo(db)
	var jobsRepo jobs.Repository
	jobsRepo = jobs.NewJobsRepo(db)

	for i := 0; i < 10; i++ {
		item := &items.Item{
			ID:  strconv.Itoa(i),
			API: sql.NullString{String: "{\"test\":\"test\"}", Valid: true},
		}

		err = itemsRepo.Save(ctx, item)
		require.NoError(t, err)
	}

	for i := 0; i < 10; i++ {
		_, err = jobsRepo.CreateJobForUnqueued(ctx, "testing")
		require.NoError(t, err)
	}

	for i := 0; i < 11; i++ {
		if i == 10 {
			job, err := jobsRepo.GetReadyJob(ctx, "testing")
			assert.Nil(t, job)
			assert.ErrorIs(t, err, jobs.ErrNoJobs)
		} else {
			job, err := jobsRepo.GetReadyJob(ctx, "testing")
			assert.NoError(t, err)
			assert.Equal(t, "running", job.Status)
			job.Finish()
			err = jobsRepo.SaveFullText(ctx, job)
			assert.NoError(t, err)
		}
	}

}
