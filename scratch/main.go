package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/lmullen/cchc/common/db"
	"github.com/lmullen/cchc/common/items"
	"github.com/lmullen/cchc/common/jobs"
)

func main() {

	ctx := context.Background()
	db, err := db.Connect(ctx, os.Getenv("CCHC_DBSTR_LOCAL"))
	if err != nil {
		log.Fatal(err)
	}
	ir := items.NewItemRepo(db)
	id := "http://www.loc.gov/item/afc1941004_sr07/"
	url := sql.NullString{
		String: "https://www.loc.gov/item/afc1941004_sr07/",
		Valid:  true,
	}

	item := &items.Item{ID: id, URL: url}

	fmt.Println(item)

	err = ir.Save(ctx, item)
	if err != nil {
		log.Fatal(err)
	}

	err = item.Fetch(http.DefaultClient)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(item)

	jr := jobs.NewJobsRepo(db)

	job := jobs.NewFullText(item.ID, "testing", false)
	fmt.Println(job)
	jr.SaveFullText(ctx, job)
	job.Start()
	jr.SaveFullText(ctx, job)
	job2, _ := jr.GetFullText(ctx, job.ID)
	fmt.Println(job2)

}
