// This program crawls the LOC.gov API for a single collection and outputs basic
//  metadata about each item.
package main

import (
	"fmt"
	"log"
	"net/http"
)

// Hard code a collection for now
const collectionSlug = "african-american-perspectives-rare-books"
const itemsPerPage = 10

func main() {

	client := &http.Client{}

	u, err := CollectionURL(collectionSlug, 1)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(u)

	collectionResults := make(chan CollectionResult, 10)

	go fetchCollectionResult(u, client, collectionResults)

	for r := range collectionResults {
		fmt.Println(r)
		for _, v := range r.Results {
			fmt.Println(v)
		}
	}

}
