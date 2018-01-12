package controller

import (
	"bytes"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/katreinhart/askify-api-v2/model"
)

// FetchAllQuestions fetch all the data from the model and handle responding
func FetchAllQuestions(w http.ResponseWriter, r *http.Request) {

	js, err := model.FetchAllQuestions()

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		if err.Error() == "Not found" {
			w.WriteHeader(http.StatusNotFound)
			w.Write(js)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something went wrong"))
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

// CreateQuestion controller function takes request and parses it, sends to model, sends response or error via JSON.
func CreateQuestion(w http.ResponseWriter, r *http.Request) {
	// Create a new buffer to read the body, then parse into a []byte
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	// Send the []byte b to the model and receive json and error
	js, err := model.CreateQuestion(b)

	// Set headers before any response is sent
	w.Header().Set("Content-Type", "application/json")

	// handle the error, if there is one; send appropriate status code
	if err != nil {
		if err.Error() == "Not found" {
			w.WriteHeader(http.StatusNotFound)
			w.Write(js)
		} else if err.Error() == "Bad request" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("{\"message\": \"Please check your inputs and try again\"}"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("{\"message\": \"Sorry, something went wrong.\"}"))
		}
	}

	// Send the success response
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

// FetchSingleQuestion fetches single question as specified by URL parameter
func FetchSingleQuestion(w http.ResponseWriter, r *http.Request) {
	// get the URL parameter from the http request
	vars := mux.Vars(r)
	id := vars["id"]

	// fetch the question and an error if there is one
	js, err := model.FetchSingleQuestion(id)

	// Set the header for the outgoing response
	w.Header().Set("Content-Type", "application/json")

	// Handle the error
	if err != nil {
		// handle it
	}

	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

// UpdateQuestion handles PUT requests to /questions/:id
func UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	// get the URL parameter from the http request
	vars := mux.Vars(r)
	id := vars["id"]

	// get the body from the request
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	js, err := model.UpdateQuestion(id, b)

	// Set the header for the outgoing response
	w.Header().Set("Content-Type", "application/json")

	// Handle the error
	if err != nil {
		// handle error
	}

	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

// FetchQueue gets all open questions in order that they were posted.
func FetchQueue(w http.ResponseWriter, r *http.Request) {
	js, err := model.FetchQueue()

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		if err.Error() == "Not found" {
			w.WriteHeader(http.StatusNotFound)
			w.Write(js)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something went wrong"))
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
