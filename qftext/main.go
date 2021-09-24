package main

import (
	"context"
	"fmt"

	"github.com/lmullen/cchc/common/items"
	log "github.com/sirupsen/logrus"
)

var app = &App{}

func main() {
	err := app.Init()
	if err != nil {
		log.WithError(err).Fatal("Error initializing application")
	}
	defer app.Shutdown()

	ir := items.NewItemRepo(app.DB)

	item, err := ir.Get(context.TODO(), "http://www.loc.gov/item/2020780885/")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(item)

	fmt.Println(item.Resources)
	fmt.Println(item.Files)

	// ir := repositories.NewItemRepo(app.DB)

	// item, err := ir.Get(context.TODO(), "http://www.loc.gov/item/scsm000087/")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(item)

	// err = GetUnqueued(context.TODO())
	// if err != nil {
	// 	log.Fatal(err)
	// }

}
