package main

import (
	"log"
	"url-shortener/db_handler"
)

type URLData struct {
	ID          string `json:"_id,omitempty" bson:"_id,omitempty"`
	URL         string `json:"url,omitempty" bson:"url,omitempty"`
	ShortCode   int    `json:"shortCode,omitempty" bson:"shortCode,omitempty"`
	CreatedAt   string `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt   string `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	AccessCount int    `json:"accessCount,omitempty" bson:"accessCount,omitempty"`
}

func main() {

	client, err := db_handler.Connect("localhost", 27017)
	if err != nil {
		log.Fatalf("Couldn't create DB client: %v", err)
	}
	defer func() {
		if err := client.Disconnect(); err != nil {
			log.Fatalf("Couldn't disconnect DB client: %v", err)
		}
	}()

	if err := client.SelectDB("urls"); err != nil {
		log.Fatalf("Error selecting db: %v", err)
	}

	collection, err := client.AddCollection("url_collection")
	if err != nil {
		log.Fatalf("Error adding collection: %v", err)
	}

}
