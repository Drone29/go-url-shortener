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
	fmt.Printf("Inserted successfully, id %s\n", id)

	new_doc := TestData{
		Name:  "test",
		Value: 1,
	}

	// list ids
	fmt.Println("Docs' IDs:")
	ids, err := collection.ListDocsIDs()
	if err != nil {
		log.Fatalf("Error obtaining list of docs: %v", err)
	}

	for _, id := range ids {
		fmt.Println(id)
	}

	err = collection.FindByID(id, &new_doc)
	if err != nil {
		log.Fatalf("Error obtaining doc from db: %v", err)
	}
	fmt.Printf("Successfully retrieved doc from db: %v", new_doc)

	err = collection.Delete(id)
	if err != nil {
		log.Fatalf("Error deleting doc %s: %v", id, err)
	}

	// list ids
	fmt.Println("Docs' IDs:")
	ids, err = collection.ListDocsIDs()
	if err != nil {
		log.Fatalf("Error obtaining list of docs: %v", err)
	}

	for _, id := range ids {
		fmt.Println(id)
	}

}
