package handlers

import (
	"backend/db"
	"backend/models"
	"backend/utils"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// HandleLogin handles user login requests
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var loginRequest models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Find user by email
	collection := db.GetDatabase().Collection("users")
	filter := bson.M{"email": loginRequest.Email}

	var user models.User
	err = collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == db.ErrNotFound {
			// User not found
			respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		// Error occurred while querying the database
		respondWithError(w, http.StatusInternalServerError, "Failed to authenticate user")
		return
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		// Passwords do not match
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Passwords match, generate JWT
	claims := models.Claims{Email: user.Email}
	token, err := utils.GenerateJWT(claims)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Send JWT as response
	respondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}
