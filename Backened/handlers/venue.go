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
	"time"
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

// HandleGetVenueByID handles requests to retrieve venue details by ID
func HandleGetVenueByID(w http.ResponseWriter, r *http.Request, venueID string) {

	if venueID == "" {
		respondWithError(w, http.StatusBadRequest, "Venue ID is required")
		return
	}

	// Fetch venue details by ID from the database
	venue, err := db.GetVenueByID(context.Background(), venueID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch venue details")
		return
	}

	// Respond with the venue details
	respondWithJSON(w, http.StatusOK, venue)
}

func HandleVenuesByCity(w http.ResponseWriter, r *http.Request) {
	// Extract the city ID from the query parameters
	cityID := r.URL.Query().Get("city")

	// If the city ID is not provided, return an error response
	if cityID == "" {
		respondWithError(w, http.StatusBadRequest, "City ID is required")
		return
	}

	// Call a function to get venues based on the city ID
	venues, err := db.GetVenuesByCity(cityID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get venues")
		return
	}

	// Respond with the list of venues
	respondWithJSON(w, http.StatusOK, venues)
}

// HandleGetOpenedVenuesByMovieID handles requests to retrieve opened venues for a specific movie
func HandleGetVenuesByMovieID(w http.ResponseWriter, r *http.Request) {
	// Parse movie ID from the request query parameters
	movieID := r.URL.Query().Get("movie_id")
	if movieID == "" {
		respondWithError(w, http.StatusBadRequest, "Movie ID is required")
		return
	}

	// Parse date from the request query parameters
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		respondWithError(w, http.StatusBadRequest, "Date is required")
		return
	}

	// Parse date string into time.Time
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid date format")
		return
	}

	// Fetch opened venues for the specified movie from the database
	openedVenues, err := db.GetOpenedVenuesByMovieID(context.Background(), movieID, date)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch opened venues")
		return
	}

	// Respond with the list of opened venues
	respondWithJSON(w, http.StatusOK, openedVenues)
}
