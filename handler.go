package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"path/filepath"
	"regexp"
	"time"
	"github.com/form3tech-oss/jwt-go"
)

var (
	mediaTypes = map[string]string{
		".jpeg": "image",
        ".jpg":  "image",
        ".gif":  "image",
        ".png":  "image",
        ".mov":  "video",
        ".mp4":  "video",
        ".avi":  "video",
        ".flv":  "video",
        ".wmv":  "video",
	}
)

var mySigningKey = []byte("ArthurPendragon")

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("An upload request received.")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")

	if r.Method == "OPTIONS" {
		return
	}

	// initiate post
	user := r.Context().Value("user")
	claims := user.(*jwt.Token).Claims
	username := claims.(jwt.MapClaims)["username"]
	p := Post{
		Username: username.(string),
		Message: r.FormValue("message"),
	}

	// get file and header
	file, header, err := r.FormFile("media_file")
	if err != nil {
		http.Error(w, "Media file is not available.", http.StatusBadRequest)
		fmt.Printf("Media file is not avaiable. %v\n", err)
		return
	}
	
	// get file type from header
	suffix := filepath.Ext(header.Filename)
	if t, ok := mediaTypes[suffix]; ok {
		p.Type = t
	} else {
		p.Type = "unknown"
	}

	// save post
	err = savePost(&p, file)
	if err != nil {
		http.Error(w, "Failed to save post to GCS or Elasticsearch.", http.StatusInternalServerError)
		fmt.Printf("Failed to save post to GCS or Elasticsearch. %v\n", err)
		return
	}
	fmt.Println("Post saved successfully.")
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("A search request received.")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")

	if r.Method == "OPTIONS" {
		return
	}

	username := r.URL.Query().Get("username")
	keywords := r.URL.Query().Get("keywords")

	var posts []Post
	var err error

	if username != "" {
		posts, err = searchPostsByUsername(username)
	} else {
		posts, err = searchPostsByKeywords(keywords)
	}
	if err != nil {
		http.Error(w, "Failed to read posts from Elasticsearch.", http.StatusInternalServerError)
		fmt.Printf("Failed to read posts from Elasticsearch. %v.\n", err)
		return
	}

	jsRes, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, "Failed to parse posts to JSON format.", http.StatusInternalServerError)
		fmt.Printf("Failed to parse posts to JSON format. %v.\n", err)
		return
	}
	w.Write(jsRes)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("A delete request received.")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")

	if r.Method == "OPTIONS" {
		return
	}

	// delete post
	id := r.URL.Query().Get("id")
	err := deletePost(id)
	if err != nil {
		http.Error(w, "Failed to delete the post by id.", http.StatusInternalServerError)
		fmt.Printf("Failed to delete the post by id. %v.\n", err)
		return
	}
	fmt.Println("Post deleted successfully.")
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("An signup request received.")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	// get user information from client
	decoder := json.NewDecoder(r.Body)
	var user User
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "Cannot decode user.", http.StatusBadRequest)
		fmt.Printf("Cannot decode user. %v\n", err)
		return
	}

	exists, correct, err := checkUser(user.Username, user.Password)
	if err != nil {
		http.Error(w, "Cannot read user from Elasticsearch.", http.StatusInternalServerError)
		fmt.Printf("Cannot read user from Elasticsearch. %v\n", err)
		return
	}
	if !exists {
		http.Error(w, "Username doesn't exist.", http.StatusBadRequest)
		fmt.Printf("Username doesn't exist.\n")
		return
	}
	if !correct {
		http.Error(w, "Wrong password.", http.StatusBadRequest)
		fmt.Printf("Wrong password.\n")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		http.Error(w, "Failed to generate token.", http.StatusInternalServerError)
		fmt.Printf("Failed to generate token. %v\n", err)
		return
	}

	w.Write([]byte(tokenString))
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("An signup request received.")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	decoder := json.NewDecoder(r.Body)
	var user User
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "Cannot decode user.", http.StatusBadRequest)
		fmt.Printf("Cannot decode user. %v\n", err)
		return
	}

	if user.Username == "" || user.Password == "" || !regexp.MustCompile(`^[a-z0-9]*$`).MatchString(user.Username) {
		http.Error(w, "Invalid username or password.", http.StatusBadRequest)
		fmt.Printf("Invalid username: %s\n", user.Username)
        return
	}

	success, err := addUser(&user)
	if err != nil {
		http.Error(w, "Cannot save user to Elasticsearch.", http.StatusInternalServerError)
		fmt.Printf("Cannot save user to Elasticsearch. %v\n", err)
		return
	}
	if !success {
		http.Error(w, "User already exists.", http.StatusBadRequest)
		fmt.Printf("User already exists.\n")
		return
	}
	fmt.Printf("User %s signed up successfully.\n", user.Username)
}