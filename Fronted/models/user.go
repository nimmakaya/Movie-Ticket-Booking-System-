package models

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Any other user-related logic

// Define a struct for login requests
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Define a struct for JWT claims
type Claims struct {
	Email string `json:"email"`
	// Add any other claims you want to include
	// (e.g., user role, expiration time, etc.)
}

// AdminUser represents the structure of an admin user document in the database
type AdminUser struct {
	Email    string `bson:"email"`
	Password string `bson:"password"`
}

// City represents the structure of a city entity
type City struct {
	Name string `json:"city_name,omitempty" bson:"city_name,omitempty"`
}

// Venue represents the structure of a venue entity
type Venue struct {
	Name            string `json:"venueName,omitempty" bson:"venue_name,omitempty"`
	CityID          string `json:"cityId,omitempty" bson:"_id,omitempty"`
	NumberOfScreens int    `json:"numberOfScreens,omitempty" bson:"number_of_screens,omitempty"`
	Address         string `json:"address,omitempty" bson:"address,omitempty"`
	ContactNumber   string `json:"contactNumber,omitempty" bson:"contact_number,omitempty"`
}
