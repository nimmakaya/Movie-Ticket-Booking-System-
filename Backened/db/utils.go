// db package

package db

import (
	"backend/models"

	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InsertCity inserts a new city record into the database
func InsertCity(ctx context.Context, city *models.City) error {
	// Get the collection
	collection := GetDatabase().Collection("cities")

	// Insert the city document into the collection
	_, err := collection.InsertOne(ctx, city)
	if err != nil {
		return err
	}

	return nil
}

// GetCityByName retrieves a city from the database by its name
func GetCityByName(ctx context.Context, cityName string) (*models.City, error) {
	var city models.City
	filter := bson.M{"city_name": cityName}

	collection := GetDatabase().Collection("cities")
	err := collection.FindOne(ctx, filter).Decode(&city)
	if err != nil {
		if err == ErrNotFound {
			// City not found
			return nil, nil
		}
		// Error occurred while querying the database
		return nil, err
	}
	return &city, nil
}

// GetCities retrieves the list of cities from the database
func GetCities() ([]bson.M, error) {
	var cities []bson.M

	// Get database collection
	collection := client.Database("movie-booking").Collection("cities")

	// Find all cities
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Decode cities from cursor
	err = cursor.All(context.Background(), &cities)
	if err != nil {
		return nil, err
	}

	return cities, nil
}

// GetCityByID retrieves a city from the database by its ID
func GetCityByID(ctx context.Context, cityID string) (*models.City, error) {
	var city models.City
	hexCityID, _ := primitive.ObjectIDFromHex(cityID)
	fmt.Println(hexCityID)
	filter := bson.M{"_id": hexCityID}
	err := GetDatabase().Collection("cities").FindOne(ctx, filter).Decode(&city)
	if err != nil {
		return nil, err
	}
	return &city, nil
}

// InsertVenue inserts a new venue into the database
func InsertVenue(ctx context.Context, venue *models.Venue) error {
	collection := GetDatabase().Collection("venues")
	_, err := collection.InsertOne(ctx, venue)
	if err != nil {
		return err
	}
	return nil
}
