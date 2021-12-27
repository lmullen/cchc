package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/lmullen/cchc/common/items"
)

func main() {

	// ctx := context.Background()
	// db, err := db.Connect(ctx, os.Getenv("CCHC_DBSTR_LOCAL"))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// ir := items.NewItemRepo(db)
	id := "http://www.loc.gov/item/afc1941004_sr07/"
	url := sql.NullString{
		String: "https://www.loc.gov/item/afc1941004_sr07/",
		Valid:  true,
	}

	item := &items.Item{ID: id, URL: url}

	fmt.Println(item)

	// err = ir.Save(ctx, item)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	err := item.Fetch(http.DefaultClient)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(item)

	texts, has := item.FullText()
	fmt.Println(texts)
	fmt.Println(has)
	fmt.Println(len(texts))

	// ir.Save(ctx, gotten)

}
