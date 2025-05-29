package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB struct holds the MongoDB client instance.
type MongoDB struct {
	Client *mongo.Client
}

// InitMongoDB connects to the MongoDB database.
func InitMongoDB(uri string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the primary to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	fmt.Println("Connected to MongoDB!")
	return &MongoDB{Client: client}, nil
}

// GetDatabase returns a specific MongoDB database instance.
// This is a method on the MongoDB struct, allowing you to easily get
// a database connection from your initialized client.
func (m *MongoDB) GetDatabase(dbName string) *mongo.Database {
	return m.Client.Database(dbName)
}
