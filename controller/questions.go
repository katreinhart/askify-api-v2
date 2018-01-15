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
	handleErrorAndRespond(js, err, w)
}

// CreateQuestion controller function takes request and parses it, sends to model, sends response or error via JSON.
func CreateQuestion(w http.ResponseWriter, r *http.Request) {
	// Create a new buffer to read the body, then parse into a []byte
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	// Send the []byte b to the model and receive json and error
	js, err := model.CreateQuestion(b)

	handleErrorAndRespond(js, err, w)
}

// FetchSingleQuestion fetches single question as specified by URL parameter
func FetchSingleQuestion(w http.ResponseWriter, r *http.Request) {
	// get the URL parameter from the http request
	vars := mux.Vars(r)
	id := vars["id"]

	// fetch the question and an error if there is one
	js, err := model.FetchSingleQuestion(id)

	handleErrorAndRespond(js, err, w)
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

	handleErrorAndRespond(js, err, w)
}

// FetchQueue gets all open questions in order that they were posted.
func FetchQueue(w http.ResponseWriter, r *http.Request) {
	js, err := model.FetchQueue()

	handleErrorAndRespond(js, err, w)
}

// FetchArchive handles request for deeply nested archive object and responds
func FetchArchive(w http.ResponseWriter, r *http.Request) {

	// Get the data from the model
	js, err := model.FetchArchive()

	handleErrorAndRespond(js, err, w)
}
