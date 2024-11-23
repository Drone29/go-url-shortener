package main

import (
	"fmt"
	"log"
	"url-shortener/db_handler"
)

func main() {

	client := db_handler.DBClient{}

	err := client.Connect("localhost", 27017)
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

}
