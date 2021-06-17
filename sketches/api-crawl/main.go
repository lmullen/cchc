// This program crawls the LOC.gov API for a single collection and outputs basic
//  metadata about each item.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Hard code a collection for now
const collectionSlug = "african-american-perspectives-rare-books"

func main() {
	u, err := CollectionURL(collectionSlug, 1)
	if err != nil {
		log.Fatalln(err)
	}

	res, err := http.Get(u)
	if err != nil {
		log.Fatalln(err)
	}

	data, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(data))

}
