package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client      *mongo.Client
	dbName      = "movie-booking"
	mongoURI    = "mongodb+srv://sushanthdats:sush2901@cluster0.qhibrct.mongodb.net/movie-booking?retryWrites=true&w=majority"
	ErrNotFound = mongo.ErrNoDocuments
)

// Init initializes the MongoDB client and connects to the database
func Init() error {
	// Set client options
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Println("Connected to MongoDB!")
	return nil
}

// GetClient returns the MongoDB client
func GetClient() *mongo.Client {
	return client
}

// GetDatabase returns the MongoDB database instance
func GetDatabase() *mongo.Database {
	return client.Database(dbName)
}
