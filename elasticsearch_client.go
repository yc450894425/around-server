package main

import (
	"fmt"
	"context"
	"github.com/olivere/elastic"
)

const (
	ES_URL = "http://10.128.0.3:9200"
)

func readFromES(query elastic.Query, index string) (*elastic.SearchResult, error) {
	// create elastic client
	client, err := elastic.NewClient(
		elastic.SetURL(ES_URL),
		elastic.SetBasicAuth("elastic", "12345678"))
	if err != nil {
		return nil, err
	}
	// search
	SearchResult, err := client.Search().
		Index(index).
		Query(query).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	return SearchResult, nil
}

func saveToES(i interface{}, index string, id string) error {
	// create elastic client
	client, err := elastic.NewClient(
		elastic.SetURL(ES_URL), 
		elastic.SetBasicAuth("elastic", "12345678"))
	if err != nil {
		return err
	}
	// index a post (using JSON serialization)
	post, err := client.Index().
		Index(index).
		Id(id).
		BodyJson(i).
		Do(context.Background())
	if err != nil {
		return err
	}
	fmt.Printf("Indexed document %s to index %s.\n", post.Id, post.Index)
	return nil
}

func deleteFromES(index string, id string) error {
	// create elastic client
	client, err := elastic.NewClient(
		elastic.SetURL(ES_URL), 
		elastic.SetBasicAuth("elastic", "12345678"))
	if err != nil {
		return err
	}
	// delete
	_, err = client.Delete().
		Index(index).
		Id(id).
		Do(context.Background())
	if err != nil {
		fmt.Println("Failed to delete the post.")
		return err
	}
	return nil
}