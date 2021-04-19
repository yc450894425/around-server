package main

import (
	"fmt"
	"reflect"
	"github.com/olivere/elastic"
)

const (
	USER_INDEX = "user"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Age int64 `json:"age"`
	Gender string `json:"gender"`
}

func checkUser(username, password string) (bool, bool, error) {
	// get the user
	searchResult, err := searchUser(username)
	if err != nil {
		return false, false, err
	}
	if searchResult.TotalHits() == 0 {
		return false, false, nil
	}

	var utype User
	for _, item := range searchResult.Each(reflect.TypeOf(utype)) {
		if u, ok := item.(User); ok {
			if u.Password == password {
				fmt.Printf("Logged in as user %s.\n", username)
				return true, true, nil
			}
		}
	}
	return true, false, nil
}

func addUser(user *User) (bool, error) {
	// de-duplication
	searchResult, err := searchUser(user.Username)
	if err != nil {
		return false, err
	}
	if searchResult.TotalHits() > 0 {
		return false, err
	}
	// add user
	err = saveToES(user, USER_INDEX, user.Username)
	if err != nil {
		return false, err
	}
	fmt.Printf("User %s is added.\n", user.Username)
	return true, nil
}

func searchUser(username string) (*elastic.SearchResult, error) {
	query := elastic.NewTermQuery("username", username)
	searchResult, err := readFromES(query, USER_INDEX)
	if err != nil {
		return nil, err
	}
	return searchResult, nil
}
