package handlers

import (
	"backend/db"
	"backend/models"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// HandleRegister handles user registration requests
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling registration request")

	// Parse request body
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println("Error decoding request body:", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate user input
	if user.Username == "" || user.Email == "" || user.Password == "" {
		log.Println("Invalid user input:", user)
		respondWithError(w, http.StatusBadRequest, "Username, email, and password are required")
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Check if the email already exists in the database
	collection := db.GetDatabase().Collection("users")
	filter := bson.M{"email": user.Email}
	var existingUser models.User
	err = collection.FindOne(context.Background(), filter).Decode(&existingUser)
	if err == nil {
		// Email already exists, inform the user
		log.Println("User with email already exists:", user.Email)
		respondWithError(w, http.StatusConflict, "Email already registered")
		return
	} else if err != nil {
		// Some error occurred while checking for existing email
		log.Println("Error checking for existing email:", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to register user")
		return
	}

	// Insert user into the database
	_, err = collection.InsertOne(context.Background(), bson.M{"username": user.Username, "email": user.Email, "password": string(hashedPassword)})
	if err != nil {
		log.Println("Error inserting user into database:", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to register user")
		return
	}

	// Respond with success message
	respondWithJSON(w, http.StatusCreated, map[string]string{"message": "User registered successfully"})
}

// respondWithError sends an error response with the given status code and message
func respondWithError(w http.ResponseWriter, status int, message string) {
	log.Println("Sending error response:", message)
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// respondWithJSON sends a JSON response with the given status code and data
func respondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	log.Println("Sending JSON response:", data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
