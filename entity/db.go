package entity

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client = nil
var cancel context.CancelFunc = nil
var ctx context.Context = nil

// GetDB returns the database
func GetDB() (*mongo.Database, context.CancelFunc) {
	if client != nil {
		return client.Database("go-test"), cancel
	}
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Failed to Ping!")
	}
	return client.Database("go-test"), cancel
}

// CloseDB closes the connection to the database
func CloseDB() {
	if client == nil {
		return
	}
	err := client.Disconnect(ctx)
	if err != nil {
		log.Printf("Could not disconnect from db: %v", err)
	}
	fmt.Println("Disconnected from DB")
}
