// This program crawls the LOC.gov API for a single collection and outputs basic
//  metadata about each item.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	res, err := client.Get(u)
	if err != nil {
		log.Fatalln(err)
	}

	data, _ := ioutil.ReadAll(res.Body)

	var result CollectionResult

	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(result)

	for _, v := range result.Results {
		fmt.Println(v)
	}

}
