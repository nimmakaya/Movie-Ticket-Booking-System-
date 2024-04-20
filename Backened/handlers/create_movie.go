package handlers

import (
	"backend/db"
	"backend/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"go.mongodb.org/mongo-driver/bson"
)

// HandleCreateMovie handles requests to create a new movie
func HandleCreateMovie(w http.ResponseWriter, r *http.Request) {
	// Parse form data including file upload
	err := r.ParseMultipartForm(10 << 20) // 10 MB max file size
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to parse form")
		return
	}

	// Get form data
	movieName := r.FormValue("name")
	releaseDate := r.FormValue("release_date")
	cast := r.FormValue("cast")
	crew := r.FormValue("crew")

	// Upload file to AWS S3
	file, handler, err := r.FormFile("poster")
	defer file.Close()
	if err != nil {
		fmt.Println("Error retrieving the file")
	} else {
		// Create a new S3 session
		sess, err := session.NewSession(&aws.Config{
			Region:      aws.String("eu-north-1"),
			Credentials: credentials.NewStaticCredentials("secretkey", "accesskey", ""),
		})

		if err != nil {
			fmt.Println("Error creating AWS session:", err)
			return
		}

		// Create S3 service client
		svc := s3.New(sess)

		// Upload file to S3 bucket
		_, err = svc.PutObject(&s3.PutObjectInput{
			Bucket: aws.String("posters-bmc"),
			Key:    aws.String(handler.Filename),
			Body:   file,
		})
		if err != nil {
			fmt.Println("Error uploading file to S3: ", err)
		} else {
			fmt.Println("File uploaded successfully")
		}
	}

	// Format release date
	releaseDateTime, err := time.Parse(time.RFC3339, releaseDate)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid release date format")
		return
	}

	// Create movie object
	movie := models.Movie{
		Name:          movieName,
		ReleaseDate:   releaseDateTime,
		ReleaseStatus: "0",
		Cast:          strings.Split(cast, ","),
		Crew:          strings.Split(crew, ","),
		PosterURL:     fmt.Sprintf("https://posters-bmc.s3.eu-north-1.amazonaws.com/%s", handler.Filename),
	}

	// Insert movie into the database
	err = db.InsertMovie(context.Background(), &movie)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create movie")
		return
	}

	// Respond with success message
	respondWithJSON(w, http.StatusCreated, map[string]string{"message": "Movie created successfully"})
}

// DefaultShowTimes represents the default show times for movies
var DefaultShowTimes = []string{"11am", "2pm", "6pm", "9pm"}

// HandleOpenMovie handles requests to open a movie
func HandleOpenMovie(w http.ResponseWriter, r *http.Request, movieID string) {
	// Check if movieID is provided
	if movieID == "" {
		respondWithError(w, http.StatusBadRequest, "Movie ID is required")
		return
	}

	// Parse the request body to get the city, venues, show times, start date, and end date
	var requestBody models.OpenMovieRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input fields
	if requestBody.CityID == "" || len(requestBody.Venues) == 0 || requestBody.StartDate == "" || requestBody.EndDate == "" {
		respondWithError(w, http.StatusBadRequest, "City ID, venues, start date, and end date are required")
		return
	}

	requestBody.ShowTimes = DefaultShowTimes

	// Parse start and end dates
	startDate, err := time.Parse(time.RFC3339, requestBody.StartDate)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid start date format")
		return
	}
	endDate, err := time.Parse(time.RFC3339, requestBody.EndDate)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid end date format")
		return
	}

	// Check venue availability
	available, err := db.CheckVenueAvailability(r.Context(), requestBody.Venues, startDate, endDate)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to check venue availability")
		return
	}
	if !available {
		respondWithError(w, http.StatusConflict, "Selected venues are not available for the specified date range")
		return
	}

	// Update the release status of the movie to 1 (opened) in the movies table
	err = db.UpdateMovieReleaseStatus(r.Context(), movieID, "1")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to open movie")
		return
	}

	// Insert the opened movie details into the opened movies table
	err = db.InsertOpenedMovie(r.Context(), movieID, requestBody.CityID, requestBody.Venues, requestBody.ShowTimes, startDate, endDate)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to insert opened movie details")
		return
	}

	// Respond with success message
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Movie opened successfully"})
}

// HandleGetUpcomingMovies handles requests to get upcoming movies with release status 0
func HandleGetUpcomingMovies(w http.ResponseWriter, r *http.Request) {
	// Fetch movies with release status 0 from the database
	movies, err := db.GetMoviesByReleaseStatus(context.Background(), "0")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch upcoming movies")
		return
	}

	// Respond with the list of upcoming movies
	respondWithJSON(w, http.StatusOK, movies)
}

// HandleGetOpenedMoviesByCity handles requests to retrieve opened movies in a specific city
func HandleGetOpenedMoviesByCity(w http.ResponseWriter, r *http.Request) {
	// Parse city ID from the request query parameters
	cityID := r.URL.Query().Get("city_id")
	if cityID == "" {
		respondWithError(w, http.StatusBadRequest, "City ID is required")
		return
	}

	// Fetch opened movies in the specified city from the database
	openedMovies, err := db.GetOpenedMoviesByCity(context.Background(), cityID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch opened movies")
		return
	}

	// Fetch complete details of each opened movie
	completeMovies := make([]bson.M, 0)
	for _, movie := range openedMovies {
		fmt.Println(movie)
		movieID := movie["movie_id"].(string)
		completeMovie, err := db.GetMovieByID(context.Background(), movieID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to fetch complete details of opened movies")
			return
		}
		completeMovies = append(completeMovies, completeMovie)
	}

	// Respond with the list of opened movies with complete details
	respondWithJSON(w, http.StatusOK, completeMovies)
}

// HandleMovieDetails handles requests to fetch details of a specific movie by ID
func HandleGetMovieDetails(w http.ResponseWriter, r *http.Request, movieID string) {
	w.Header().Set("Content-Type", "application/json")

	var movie models.Movie

	// Fetch movie details from the database
	err := db.GetMovieDetailsFromDB(movieID, &movie)
	if err != nil {
		if err == db.ErrNotFound {
			// If movie not found, return 404 Not Found
			w.WriteHeader(http.StatusNotFound)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch movie")
		return
	}

	respondWithJSON(w, http.StatusOK, movie)
}

// HandleMovieDetails handles requests to fetch details of a specific movie by movie name
func GetMovieByName(w http.ResponseWriter, r *http.Request, movieName string) {
	w.Header().Set("Content-Type", "application/json")

	var movie models.Movie

	// Fetch movie details from the database
	err := db.GetMovieDetailsByName(movieName, &movie)
	if err != nil {
		if err == db.ErrNotFound {
			// If movie not found, return 404 Not Found
			w.WriteHeader(http.StatusNotFound)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch movie")
		return
	}

	respondWithJSON(w, http.StatusOK, movie)
}
