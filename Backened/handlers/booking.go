package handlers

import (
	"backend/db"     // Import your database package
	"backend/models" // Import your models package
	"encoding/json"
	"fmt"
	"net/http"
)

// CreateBooking handles the creation of a new booking
func CreateBooking(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse request body into booking struct
	var booking models.Booking
	err := json.NewDecoder(r.Body).Decode(&booking)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode request body")
		return
	}

	// Insert booking into the database
	err = db.InsertBooking(r.Context(), &booking)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to insert booking into database")
		return
	}

	// Return success response
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Booking created successfully",
		"booking": booking,
	})
}

// GetBookedSeats handles the retrieval of booked seats for a given venue ID and showtime
func GetBookedSeats(w http.ResponseWriter, r *http.Request, venueID, showTime, date string) {
	w.Header().Set("Content-Type", "application/json")

	// Validate input parameters
	if venueID == "" || showTime == "" || date == "" {
		http.Error(w, "Venue ID and showtime are required", http.StatusBadRequest)
		return
	}

	// Retrieve booked seats from the database for the given venue ID and showtime
	bookedSeats, err := db.GetBookedSeatsDB(r.Context(), venueID, showTime, date)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get bookings from database")
		return
	}

	fmt.Println(bookedSeats)

	// Return booked seats as JSON response
	json.NewEncoder(w).Encode(bookedSeats)
}

// GetBookingsByUser handles the retrieval of bookings based on the 'user' field
func GetBookingsByUser(w http.ResponseWriter, r *http.Request) {

	var request struct {
		User string `json:"user"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		// Handle decoding error
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")

	// Retrieve bookings from the database for the given user
	bookings, err := db.GetBookingsByUser(r.Context(), request.User)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get bookings from database")
		return
	}

	// Return bookings as JSON response
	json.NewEncoder(w).Encode(bookings)
}

// DeleteBooking handles the deletion of a booking
func DeleteBooking(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse request body into booking struct
	var cancelReq models.CancelBookingRequest
	err := json.NewDecoder(r.Body).Decode(&cancelReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode request body")
		return
	}

	fmt.Println(cancelReq)

	// Delete booking from the database
	err = db.DeleteBooking(r.Context(), &cancelReq)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete booking from database")
		return
	}

	// Return success response
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Booking deleted successfully",
		"booking": cancelReq,
	})
}

/*
// GetBookings handles the retrieval of all bookings
func () GetBookings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Retrieve all bookings from the database
	bookings, err := db.GetAllBookings(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve bookings from database: %v", err), http.StatusInternalServerError)
		return
	}

	// Return bookings as JSON response
	json.NewEncoder(w).Encode(bookings)
}

// GetBookingByID handles the retrieval of a booking by ID
func (bh *BookingHandler) GetBookingByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get booking ID from URL parameters
	vars := mux.Vars(r)
	bookingID := vars["id"]

	// Retrieve booking from the database by ID
	booking, err := db.GetBookingByID(r.Context(), bookingID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve booking from database: %v", err), http.StatusInternalServerError)
		return
	}

	// Check if booking is nil (not found)
	if booking == nil {
		http.Error(w, "Booking not found", http.StatusNotFound)
		return
	}

	// Return booking as JSON response
	json.NewEncoder(w).Encode(booking)
}
*/
