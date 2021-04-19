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

func main() {
	client, err := elastic.NewClient(
		elastic.SetURL(ES_URL),
		elastic.SetBasicAuth("elastic", "12345678"))
	if err != nil {
		panic(err)
	}

	// create post index
	exists, err := client.IndexExists(POST_INDEX).Do(context.Background())
	if err != nil {
		panic(err)
	}
	if !exists {
		mapping := `{
			"mappings":{
				"properties":{
					"id": {"type": "keyword" },
					"username": { "type": "keyword" },
					"message": { "type": "text" },
					"url": { "type": "keyword", "index": false },
					"type": { "type": "keyword", "index": false }
				}
			}
		}`
		_, err := client.CreateIndex(POST_INDEX).Body(mapping).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}

	// create user index
	exists, err = client.IndexExists(USER_INDEX).Do(context.Background())
	if err != nil {
		panic(err)
	}
	if !exists {
		mapping := `{
			"mappings":{
				"properties":{
					"username": { "type": "keyword" },
					"password": { "type": "keyword", "index": false },
					"age": { "type": "long", "index": false },
					"gender": { "type": "keyword", "index": false }
				}
			}
		}`
		_, err := client.CreateIndex(USER_INDEX).Body(mapping).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Indices are created.")
}