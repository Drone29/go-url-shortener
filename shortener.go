package main

import (
	"fmt"
	"log"
	"url-shortener/db_handler"
	"url-shortener/http_handler"
	"url-shortener/url_data"
)

type URLData = url_data.URLData

func main() {

	fmt.Println("Connecting to db...")
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

	collection, err := client.GetCollection("url_collection")
	if err != nil {
		log.Fatalf("Error adding collection: %v", err)
	}

	fmt.Println("Listening on port 8080...")

	http_handler.Start(8080, collection)

	fmt.Println("Stopped listening, disconnecting from db")

	client.Disconnect()

}
