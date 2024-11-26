package http_handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"url-shortener/db_handler"
	"url-shortener/url_data"
	"url-shortener/url_generator"
)

type URLData = url_data.URLData
type DB = db_handler.DBCollection

const shortURLLen int = 6

var db *DB

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

// register new url
func handlePOST(w http.ResponseWriter, r *http.Request) {

	// handle panic
	defer func() {
		if r := recover(); r != nil {
			http.Error(w, fmt.Sprintf("Internal error:\n%v", r), http.StatusInternalServerError) //500
		}
	}()

	switch r.URL.Path {
	case "/shorten", "/shorten/":
		// read body
		body := readBody(r)

		record := URLData{}
		// convert body to json
		err := json.Unmarshal(body, &record)
		if err != nil {
			// return http.StatusBadRequest in case of validation errors
			http.Error(w, "Invalid request", http.StatusBadRequest) //400
			return
		}
		// check if such record already exists
		// use parsed record as both filter and result
		if err = db.FindOne(record, &record); err != nil {
			panic(fmt.Sprintf("Error accessing db:\n%v", err))
		}
		// already exists
		if record.ID != "" {
			http.Error(w, "URL already exists", http.StatusNotModified) //304
			return
		}
		// set missing properties
		record.CreatedAt = time.Now()
		record.UpdatedAt = record.CreatedAt
		record.ShortCode = url_generator.GenerateShortURL(shortURLLen)
		// store new record in the db
		record.ID, err = db.InsertOne(record)
		if err != nil {
			panic(fmt.Sprintf("Error inserting into db:\n%v", err))
		}
		// return response
		w.WriteHeader(http.StatusCreated) //201
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprintf(w, "%s", record) // return new record as JSON
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "%s not found\n", r.URL.Path)
	}
}

// obtain registered url
func handleGET(w http.ResponseWriter, r *http.Request) {
	tokens := tokenizePath(r.URL.Path)
	switch len(tokens) {
	case 2:
		short_url := tokens[1]
		// TODO: retrieve short url from db
		// TODO: return 404 if short url not found
		w.WriteHeader(http.StatusOK) //200
		fmt.Fprintf(w, "Retrieved short url from db %s\n", short_url)
	case 3:
		if tokens[2] == "stats" {
			handleStats(tokens[1], w, r) // stats
		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "%s not found\n", r.URL.Path)
		}
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s not found\n", r.URL.Path)
	}
}

// update registered url
func handlePUT(w http.ResponseWriter, r *http.Request) {
	tokens := tokenizePath(r.URL.Path)
	switch len(tokens) {
	case 2:
		short_url := tokens[1]
		// TODO: update short url in db
		// TODO: return 404 if short url not found
		// TODO: return 400 bad request in case of validation errors
		w.WriteHeader(http.StatusOK) //200
		fmt.Fprintf(w, "Updated short url %s\n", short_url)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s not found\n", r.URL.Path)
	}
}

// remove registered url
func handleDELETE(w http.ResponseWriter, r *http.Request) {
	tokens := tokenizePath(r.URL.Path)
	switch len(tokens) {
	case 2:
		short_url := tokens[1]
		// TODO: delete short url from db
		// TODO: return 404 if short url not found
		w.WriteHeader(http.StatusNoContent) //204
		fmt.Fprintf(w, "Deleted short url %s\n", short_url)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s not found\n", r.URL.Path)
	}
}

// get statistics
func handleStats(short_url string, w http.ResponseWriter, r *http.Request) {
	//TODO: retrieve url from db and return its info
	w.WriteHeader(http.StatusOK) //200
	fmt.Fprintf(w, "Stats for url %s\n", short_url)
}

// handle http requests
func shorten(w http.ResponseWriter, r *http.Request) {
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
func Start(port int, collection *DB) error {
	if collection == nil {
		return fmt.Errorf("db collection is nil")
	}
	db = collection
	http.HandleFunc("/shorten", shorten)
	http.HandleFunc("/shorten/", shorten)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	return nil
}
