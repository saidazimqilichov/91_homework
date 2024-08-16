package repository

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)
var client *mongo.Client
func NewMngoClient(mongoURI string) (*mongo.Client, error) {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("MongoDB connection error:", err)
		return nil, err
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal("MongoDB ping error:", err)
		return nil, err
	}

	log.Println("MongoDB connected successfully")
	return client, nil
}
