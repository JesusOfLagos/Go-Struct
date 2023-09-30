package main

import (
	"context"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

// Post represents a post document in MongoDB.
type Post struct {
	ID      string `json:"id" bson:"_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}




func createPost(w http.ResponseWriter, r *http.Request) {
	// Parse request body into a Post object
	var post Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert the post into MongoDB
	postCollection := mongoClient.Database("mydb").Collection("posts")
	_, err = postCollection.InsertOne(context.Background(), post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func getPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]

	// Find the post by ID in MongoDB
	postCollection := mongoClient.Database("mydb").Collection("posts")
	var post Post
	err := postCollection.FindOne(context.Background(), bson.M{"_id": postID}).Decode(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Serialize the post as JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}