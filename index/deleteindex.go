package main

import (
	"context"
	"fmt"

	"github.com/olivere/elastic"
)

const (
	POST_INDEX = "post"
	USER_INDEX = "user"
	ES_URL = "http://10.128.0.3:9200"
)

func deleteIndex(index string) error {
	// create Elasticsearch client
	client, err := elastic.NewClient(
		elastic.SetURL(ES_URL),
		elastic.SetBasicAuth("elastic", "12345678"))
	if err != nil {
		panic(err)
	}

	// delete index
	_, err = client.DeleteIndex(index).Do(context.Background())
	if err != nil {
		panic(err)
	}
	return nil
}

func main() {
	err := deleteIndex(POST_INDEX)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Index %s is deleted.\n", POST_INDEX)
}

