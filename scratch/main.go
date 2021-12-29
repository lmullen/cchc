package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/lmullen/cchc/common/db"
	"github.com/lmullen/cchc/common/items"
)

func main() {
	rc := retryablehttp.NewClient()
	rc.RetryWaitMin = 10 * time.Second
	rc.RetryWaitMax = 2 * time.Minute
	rc.RetryMax = 6
	rc.HTTPClient.Timeout = 60 * time.Second
	rc.Logger = nil
	client := rc.StandardClient()

	ctx := context.Background()
	db, err := db.Connect(ctx, os.Getenv("CCHC_DBSTR_LOCAL"), "cchc-scratch")
	if err != nil {
		log.Fatal(err)
	}
	ir := items.NewItemRepo(db)
	id := "http://www.loc.gov/item/2020771054/"
	url := sql.NullString{
		String: "https://www.loc.gov/item/2020771054/",
		Valid:  true,
	}

	test := &items.Item{ID: id, URL: url}

	fmt.Println(test)

	err = ir.Save(ctx, test)
	if err != nil {
		log.Fatal(err)
	}

	gotten, err := ir.Get(ctx, id)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(gotten)

	fmt.Println(gotten.Fetched())

	err = gotten.Fetch(client)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(gotten)

	ir.Save(ctx, gotten)

}
