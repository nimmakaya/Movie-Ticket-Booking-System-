package handlers

import (
	"backend/db"
	"backend/models"
	"backend/utils"
	"context"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

// handleAdminLogin handles the authentication process for admin users
func HandleAdminLogin(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var loginRequest models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Find admin user by email
	collection := db.GetDatabase().Collection("admin")
	filter := bson.M{"email": loginRequest.Email}
	var adminUser models.AdminUser
	err = collection.FindOne(context.Background(), filter).Decode(&adminUser)
	if err != nil {
		if err == db.ErrNotFound {
			// Admin user not found
			respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		// Error occurred while querying the database
		respondWithError(w, http.StatusInternalServerError, "Failed to authenticate admin user")
		return
	}

	// Compare passwords
	if loginRequest.Password != adminUser.Password {
		// Passwords do not match
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Passwords match, generate JWT
	claims := models.Claims{Email: adminUser.Email}
	token, err := utils.GenerateJWT(claims)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Send JWT as response
	respondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}
