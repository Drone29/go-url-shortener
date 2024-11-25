package http_handler

import (
	"fmt"
	"net/http"
	"strings"
	"url-shortener/url_data"
	"url-shortener/url_generator"
)

type URLData = url_data.URLData

// helpers
func tokenizePath(path string) []string {
	path = strings.Trim(path, "/") // trim leading and trailing /s
	return strings.Split(path, "/")
}

// register new url
func handlePOST(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/shorten", "/shorten/":
		// TODO: store new record in the db
		short_url := url_generator.GenerateShortURL(6)
		// TODO: return http.StatusBadRequest in case of validation errors
		w.WriteHeader(http.StatusCreated) //201
		fmt.Fprintf(w, "Added new url %s to db\n", short_url)
	default:
		w.WriteHeader(http.StatusNotFound)
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

func Start(port int) {
	http.HandleFunc("/shorten", shorten)
	http.HandleFunc("/shorten/", shorten)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
