package http_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"url-shortener/db_interface"
	"url-shortener/url_data"
	"url-shortener/url_generator"
)

type URLData = url_data.URLData
type DB = db_interface.IDBCollection

const shortURLLen int = 6

var db DB
var server *http.Server

// helpers
func tokenizePath(path string) []string {
	path = strings.Trim(path, "/") // trim leading and trailing /s
	return strings.Split(path, "/")
}

func readBody(r *http.Request) []byte {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(fmt.Sprintf("Error reading body:\n%v", err))
	}
	defer r.Body.Close()
	return body
}

func recordFromBody(r *http.Request) URLData {
	// read body
	body := readBody(r)
	log.Printf("[DEBUG] Request %s", string(body))
	record := URLData{}
	// convert body to json
	err := json.Unmarshal(body, &record)
	if err != nil {
		panic(httpErr{code: http.StatusBadRequest, descr: fmt.Sprintf("Error processing request: %v", err)}) //400
	}
	return record
}

func sendJsonResponse(w http.ResponseWriter, status int, record URLData) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, "%s", record)
	log.Printf("[DEBUG] Response %s", record)
}

func handleDBErrors(err error) {
	if err != nil {
		if err == db_interface.ErrNoDocuments {
			panic(httpErr{code: http.StatusNotFound, descr: "No records found"})
		} else {
			panic(httpErr{code: http.StatusInternalServerError, descr: fmt.Sprintf("DB error: %v", err)})
		}
	}
}

// register new url
func handlePOST(w http.ResponseWriter, r *http.Request) {

	switch r.URL.Path {
	case "/shorten", "/shorten/":
		record := recordFromBody(r)
		// check if such record already exists
		// use parsed record as both filter and result
		log.Printf("[DEBUG] Looking for record in db...")
		err := db.FindOne(record, &record)

		switch err {
		case nil:
			log.Printf("[DEBUG] Record already exists")
			sendJsonResponse(w, http.StatusOK, record) //200
			return
		case db_interface.ErrNoDocuments:
			log.Printf("[DEBUG] Record doesn't exists, proceeding")
		default:
			panic(httpErr{code: http.StatusInternalServerError, descr: fmt.Sprintf("DB error: %v", err)}) // 500
		}

		// set missing properties
		record.CreatedAt = time.Now()
		record.UpdatedAt = record.CreatedAt
		record.ShortCode = url_generator.GenerateShortURL(shortURLLen)
		record.AccessCount = 0
		// store new record in the db
		log.Printf("[DEBUG] Inserting record into db...")
		record.ID, err = db.InsertOne(record)
		handleDBErrors(err)
		// return response
		sendJsonResponse(w, http.StatusCreated, record) //201
	default:
		http.Error(w, fmt.Sprintf("Not found %s", r.URL.Path), http.StatusNotFound) //404
	}
}

// get statistics
func retrieveRecord(short_url string, w http.ResponseWriter, include_ac bool) {
	record := URLData{
		ShortCode: short_url,
	}
	// retrieve short url from db
	log.Printf("[DEBUG] Looking for record in db...")
	err := db.FindOne(record, &record)
	record.IncludeAccessCountInJSON(include_ac)
	handleDBErrors(err)
	// if not stats request, update count
	if !include_ac {
		new_rec := record
		new_rec.AccessCount++
		err = db.UpdateOne(record, new_rec)
		handleDBErrors(err)
	}
	sendJsonResponse(w, http.StatusOK, record) // 200
}

// obtain registered url
func handleGET(w http.ResponseWriter, r *http.Request) {

	tokens := tokenizePath(r.URL.Path)
	switch len(tokens) {
	case 2:
		retrieveRecord(tokens[1], w, false)
	case 3:
		if tokens[2] == "stats" {
			retrieveRecord(tokens[1], w, true) // stats
		} else {
			http.Error(w, fmt.Sprintf("Not found %s", r.URL.Path), http.StatusNotFound) //404
		}
	default:
		http.Error(w, fmt.Sprintf("Not found %s", r.URL.Path), http.StatusNotFound) //404
	}
}

// update registered url
func handlePUT(w http.ResponseWriter, r *http.Request) {
	tokens := tokenizePath(r.URL.Path)
	switch len(tokens) {
	case 2:
		replaceWhat := URLData{
			ShortCode: tokens[1],
		}
		replaceWith := recordFromBody(r)
		replaceWith.ShortCode = replaceWhat.ShortCode
		replaceWith.UpdatedAt = time.Now()
		replaceWith.AccessCount = replaceWhat.AccessCount + 1
		err := db.UpdateOne(replaceWhat, replaceWith)
		handleDBErrors(err)
		sendJsonResponse(w, http.StatusOK, replaceWith) // 200
	default:
		http.Error(w, fmt.Sprintf("Not found %s", r.URL.Path), http.StatusNotFound) //404
	}
}

// remove registered url
func handleDELETE(w http.ResponseWriter, r *http.Request) {
	tokens := tokenizePath(r.URL.Path)
	switch len(tokens) {
	case 2:
		short_url := tokens[1]
		record := URLData{
			ShortCode: short_url,
		}
		err := db.DeleteOne(record)
		handleDBErrors(err)
		w.WriteHeader(http.StatusNoContent) //204
		fmt.Fprintf(w, "Deleted %s\n", short_url)
	default:
		http.Error(w, fmt.Sprintf("Not found %s", r.URL.Path), http.StatusNotFound) //404
	}
}

// recover function
func recover_hdl(w http.ResponseWriter) {
	if r := recover(); r != nil {
		log.Printf("[ERROR] %v", r)
		switch err := r.(type) {
		case httpErr:
			http.Error(w, err.descr, err.code)
		default:
			http.Error(w, fmt.Sprintf("Internal error: %v", err), http.StatusInternalServerError) //500
		}
	}
}

// handle http requests
func shorten(w http.ResponseWriter, r *http.Request) {
	// handle panic
	defer recover_hdl(w)

	switch r.Method {
	case "POST":
		handlePOST(w, r)
	case "GET":
		handleGET(w, r)
	case "PUT":
		handlePUT(w, r)
	case "DELETE":
		handleDELETE(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// start server
func Start(port int, collection DB) {
	if collection == nil {
		log.Fatalf("[ERROR] db collection is nil")
	}
	db = collection
	mux := http.NewServeMux()
	// Register handler functions with the ServeMux
	mux.HandleFunc("/shorten", shorten)
	mux.HandleFunc("/shorten/", shorten)

	server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Printf("[ERROR] %v", err)
	}
}

// shutdown server
func ShutDown() {
	log.Println("[DEBUG] Shutting down gracefully")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("[ERROR] %v", err)
	}
	log.Println("[DEBUG] Server shut down")
}
