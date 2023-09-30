package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cloudinary/cloudinary-go"
	"github.com/gorilla/mux"
)

var cld *cloudinary.Cloudinary



func main() {
	// Initialize MongoDB client
	mongoClientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	mongoClient, err = mongo.Connect(context.Background(), mongoClientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer mongoClient.Disconnect(context.Background())

	// Initialize Cloudinary client
	cld, err = cloudinary.NewFromParams("cloud_name", "api_key", "api_secret")
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the router
	router := mux.NewRouter()

	// Define API routes
	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/posts", createPost).Methods("POST")
	router.HandleFunc("/posts/{id}", getPost).Methods("GET")
	router.HandleFunc("/images", uploadImage).Methods("POST")

	// Start the server
	http.Handle("/", router)
	fmt.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}