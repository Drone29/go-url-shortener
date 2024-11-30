package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"url-shortener/backend"
	"url-shortener/db_handler"
	"url-shortener/url_data"
)

type URLData = url_data.URLData

func main() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()

	fmt.Println("Connecting to db...")
	client, err := db_handler.Connect("localhost", 27017)
	if err != nil {
		panic(err)
	}
	// disconnect db upon exit
	db_disconnect := func() {
		fmt.Println("Disconnecting DB...")
		if err = client.Disconnect(); err != nil {
			log.Printf("Couldn't disconnect DB client: %v", err)
		}
	}
	defer db_disconnect()

	if err := client.SelectDB("urls"); err != nil {
		panic(err)
	}

	collection, err := client.GetCollection("url_collection")
	if err != nil {
		panic(err)
	}

	fmt.Println("Listening on port 8080...")

	go backend.StartBackend(8080, collection)

	// add signal handler
	quit := make(chan os.Signal, 1)                    // create a channel for signals
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM) // relay SIGINT, SIGTERM signals to quit channel
	// wait for signal
	<-quit
	backend.ShutDownBackend()
}
