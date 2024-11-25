package http_handler

import (
	"fmt"
	"net/http"
	"strings"
	"url-shortener/url_data"
)

type URLData = url_data.URLData

func handlePOST(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/shorten", "/shorten/":
		//TODO: store new record in the db

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Added new url to db\n")
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func handleGET(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	path = strings.Trim(path, "/")
	tokens := strings.Split(path, "/")
	if len(tokens) == 2 {
		short_url := tokens[1]
		// TODO: retrieve short url from db
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Retrieved short url from db %s\n", short_url)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method not allowed for this path\n")
	}
}

func handlePUT(w http.ResponseWriter, r *http.Request) {

}

func handleDELETE(w http.ResponseWriter, r *http.Request) {

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
