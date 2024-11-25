package main

import (
	"encoding/json"
	"fmt"
	"log"
	"url-shortener/db_handler"
	"url-shortener/http_handler"
	"url-shortener/url_data"
)

type URLData = url_data.URLData

func main() {

	ud := URLData{
		ID:        "12",
		URL:       "https://12443",
		ShortCode: "12345",
	}
	ud.IncludeAccessCountInJSON(true)
	var ud_n URLData
	fmt.Println(ud)
	json.Unmarshal([]byte(`{"_id":"44","createdAt":"", "updatedAt":""}`), &ud_n)

	fmt.Println(ud_n)

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

	http_handler.Start(8080)

	// collection, err := client.AddCollection("url_collection")
	// if err != nil {
	// 	log.Fatalf("Error adding collection: %v", err)
	// }

}
