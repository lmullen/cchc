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

	page1, err := fetchCollectionResult(u, client)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(page1)

	for _, v := range page1.Results {
		fmt.Println(v)
	}

}
