package main

import (
	"fmt"
	"log"
	"url-shortener/db_handler"
)

type TestData struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
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

	// get db names
	dbs, err := client.GetDBNames()
	if err != nil {
		log.Fatalf("Couldn't obtain db names: %v", err)
	}
	for _, db := range dbs {
		fmt.Println(db)
	}

	if err := client.SelectDB("urls"); err != nil {
		log.Fatalf("Error selecting db: %v", err)
	}

	collection, err := client.AddCollection("url_collection")
	if err != nil {
		log.Fatalf("Error adding collection: %v", err)
	}

	// Test
	doc := TestData{
		Name:  "test",
		Value: 1,
	}

	id, err := collection.InsertOne(doc)
	if err != nil {
		log.Fatalf("Error inserting doc into collection: %v", err)
	}
	fmt.Printf("Inserted successfully, id %s", id)

	new_doc, err := db_handler.FindByID[TestData](collection, id)
	if err != nil {
		log.Fatalf("Error obtaining doc from db: %v", err)
	}
	fmt.Printf("Successfully retrieved doc from db: %v", new_doc)

}
