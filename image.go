package main

import (
	"encoding/json"
	"net/http"

	"github.com/cloudinary/cloudinary-go/api/admin"
)

// Image represents an image uploaded to Cloudinary.
type Image struct {
	ID        string `json:"id"`
	PublicID  string `json:"public_id"`
	URL       string `json:"url"`
	Thumbnail string `json:"thumbnail"`
}

// Add your image-related API routes here

func uploadImage(w http.ResponseWriter, r *http.Request) {
	// Parse image file from request
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Upload image to Cloudinary
	uploadResult, err := cld.Upload.Upload(file, admin.UploadParams{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create an Image object and store it in MongoDB
	image := Image{
		ID:        uploadResult.PublicID,
		PublicID:  uploadResult.PublicID,
		URL:       uploadResult.URL,
		Thumbnail: uploadResult.SecureURL,
	}

	imageCollection := mongoClient.Database("mydb").Collection("images")
	_, err = imageCollection.InsertOne(context.Background(), image)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Serialize the uploaded image as JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(image)
}


func deleteImage (w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageID := vars["id"]

	// Delete the image from Cloudinary
	_, err := cld.Upload.Destroy(imageID, admin.DestroyParams{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete the image from MongoDB
	imageCollection := mongoClient.Database("mydb").Collection("images")
	_, err = imageCollection.DeleteOne(context.Background(), bson.M{"_id": imageID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageID := vars["id"]

	// Find the image by ID in MongoDB
	imageCollection := mongoClient.Database("mydb").Collection("images")
	var image Image
	err := imageCollection.FindOne(context.Background(), bson.M{"_id": imageID}).Decode(&image)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Serialize the image as JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(image)
}


func getImages(w http.ResponseWriter, r *http.Request) {
	// Find all images in MongoDB
	imageCollection := mongoClient.Database("mydb").Collection("images")
	cursor, err := imageCollection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	// Serialize the images as JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cursor)
}


