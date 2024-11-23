package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Couldn't create MongoDB client: %v", err)
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatalf("Couldn't disconnect MongoDB client: %v", err)
		}
	}()

	// access db and collection
	dbs, _ := client.ListDatabaseNames(ctx, bson.D{})

	for _, db := range dbs {
		fmt.Println(db)
	}

}
