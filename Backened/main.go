package main

import (
	"backend/db"
	"backend/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {

	r := mux.NewRouter()

	// Define your routes here
	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleRegister(w, r)
	}).Methods("POST")

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleLogin(w, r)
	}).Methods("POST")

	r.HandleFunc("/admin-login", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleAdminLogin(w, r)
	}).Methods("POST")

	r.HandleFunc("/cities", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleCreateCity(w, r)
	}).Methods("POST")

	r.HandleFunc("/cities", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleGetCities(w, r)
	}).Methods("GET")

	r.HandleFunc("/venues", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleCreateVenue(w, r)
	}).Methods("POST")

	// Enable CORS with default options
	handler := cors.Default().Handler(r)

	db.Init()

	// Start the server
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
