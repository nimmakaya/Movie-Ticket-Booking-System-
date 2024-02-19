package handlers

import (
	"backend/db"
	"backend/models"
	"context"
	"encoding/json"

	"net/http"
)

// HandleCreateCity handles requests to create a new city
func HandleCreateCity(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var city models.City
	err := json.NewDecoder(r.Body).Decode(&city)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Check if the city already exists
	existingCity, err := db.GetCityByName(context.Background(), city.Name)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to check for existing city")
		return
	}
	if existingCity != nil {
		respondWithError(w, http.StatusConflict, "City already exists")
		return
	}

	// Insert city into the database
	err = db.InsertCity(context.Background(), &city)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create city")
		return
	}

	// Respond with success message
	respondWithJSON(w, http.StatusCreated, map[string]string{"message": "City created successfully"})
}

// HandleGetCities handles requests to get the list of cities
func HandleGetCities(w http.ResponseWriter, r *http.Request) {
	cities, err := db.GetCities()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get cities")
		return
	}

	// Respond with the list of cities
	respondWithJSON(w, http.StatusOK, cities)
}
