package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoDB(uri, database string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	db := client.Database(database)
	log.Println("Connected to MongoDB")
	return &MongoDB{Client: client, Database: db}, nil
}

func InitMongoDB(uri string) (*MongoDB, error) {
	return NewMongoDB(uri, "live_sports")
}

func (m *MongoDB) Disconnect(ctx context.Context) error {
	if m.Client != nil {
		return m.Client.Disconnect(ctx)
	}
	return nil
}

func (m *MongoDB) GetDatabase() *mongo.Database {
	return m.Database
}