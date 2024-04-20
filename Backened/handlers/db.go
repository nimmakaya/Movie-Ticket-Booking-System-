package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"backend/db"
	"go.mongodb.org/mongo-driver/bson"
)

// HandleAddPosterURL handles requests to add the poster_url field to movie documents
func HandleAddPosterURL(w http.ResponseWriter, r *http.Request) {
	// Connect to MongoDB
	client := db.GetClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get a handle to the movies collection
	collection := client.Database("movie-booking").Collection("movies")

	// Define an update operation to add the poster_url field
	update := bson.M{"$set": bson.M{"poster_url": nil}} // You can set a default value here if needed

	// Update all documents in the collection to add the poster_url field
	result, err := collection.UpdateMany(ctx, bson.M{}, update)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add poster_url field: %s", err), http.StatusInternalServerError)
		return
	}

	// Respond with the number of documents updated
	response := struct {
		Message string `json:"message"`
		Updated int64  `json:"updated"`
	}{
		Message: "Successfully added poster_url field to movie documents",
		Updated: result.ModifiedCount,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
