package http_handler

import (
	"fmt"
	"net/http"
)

func handlePOST(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/shorten", "/shorten/":
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Added new url to db\n")
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func handleGET(w http.ResponseWriter, r *http.Request) {

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
