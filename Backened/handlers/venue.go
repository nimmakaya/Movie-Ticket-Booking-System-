package handlers

import (
	"backend/db"
	"backend/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// HandleCreateVenue handles requests to create a new venue
func HandleCreateVenue(w http.ResponseWriter, r *http.Request) {

	// Read the request body into a buffer
	bodyBuffer, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to read request body")
		return
	}

	// Print the raw request body
	fmt.Println(string(bodyBuffer))

	// Restore the request body for decoding
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBuffer))

	// Parse request body
	var venue models.Venue
	err = json.NewDecoder(r.Body).Decode(&venue)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	fmt.Println(venue)

	// Check if the city exists
	_, err = db.GetCityByID(context.Background(), venue.CityID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid city ID")
		return
	}

	// Insert venue into the database
	err = db.InsertVenue(context.Background(), &venue)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create venue")
		return
	}

	// Respond with success message
	respondWithJSON(w, http.StatusCreated, map[string]string{"message": "Venue created successfully"})
}
