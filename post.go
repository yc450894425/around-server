package main

import (
	"fmt"
	"mime/multipart"
	"reflect"
	"github.com/olivere/elastic"
	"github.com/pborman/uuid"
)

const (
	POST_INDEX = "post"
)

type Post struct {
	Id string `json:"id"`
	Username string `json:"username"`
	Message string `json:"message"`
	Url string `json:"url"`
	Type string `json:"type"`
}

func searchPostsByUsername(username string) ([]Post, error) {
	query := elastic.NewTermQuery("username", username)
	searchResult, err := readFromES(query, POST_INDEX)
	if err != nil {
		return nil, err
	}
	
	return getPostFromSearchResult(searchResult), nil
}

func searchPostsByKeywords(keywords string) ([]Post, error) {
	query := elastic.NewMatchQuery("message", keywords)
	query.Operator("AND")
	if keywords == "" {
		query.ZeroTermsQuery("all")
	}
	searchResult, err := readFromES(query, POST_INDEX)
	if err != nil {
		return nil, err
	}

	return getPostFromSearchResult(searchResult), nil
}

func getPostFromSearchResult(searchResult *elastic.SearchResult) []Post {
	var ptype Post
	var posts []Post

	for _, item := range searchResult.Each(reflect.TypeOf(ptype)) {
		if p, ok := item.(Post); ok {
			posts = append(posts, p)
		}
	}

	return posts
}

func savePost(post *Post, file multipart.File) error {
	// generate uuid for the post and file
	id := uuid.New()
	post.Id = id
	// save media file to GCS
	medialink, err := saveToGCS(file, id)
	if err != nil {
		return err
	}
	post.Url = medialink
	// save post to Elasticsearch
	err = saveToES(post, POST_INDEX, id)
	if err != nil {
		return err
	}
	fmt.Printf("Post and file saved to ES and GCS: %s\n", post.Message)
	return nil
}

func deletePost(id string) error {
	// get the post
	query := elastic.NewTermQuery("id", id)
	searchResult, err := readFromES(query, POST_INDEX)
	if err != nil {
		return err
	}
	if searchResult.TotalHits() == 0 {
		return nil
	}
	post := getPostFromSearchResult(searchResult)[0]
	// delete from Elasticsearch
	err = deleteFromES(POST_INDEX, id)
	if err != nil {
		return err
	}
	// delete from GCS
	if post.Url == "" {
		return nil
	}
	err = deleteFromGCS(id)
	if err != nil {
		return err
	}
	return nil
}