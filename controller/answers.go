package controller

import (
	"bytes"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/katreinhart/askify-api-v1/model"
)

// FetchQuestionAnswers fetches all the answers to a given question
func FetchQuestionAnswers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	qid := vars["id"]

	js, err := model.FetchQuestionAnswers(qid)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		// handle error
	}
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

// CreateAnswer posts an answer to a question
func CreateAnswer(w http.ResponseWriter, r *http.Request) {
	// get the data from the request body
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	js, err := model.CreateAnswer(b)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		// handle error
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(js)
}

// FetchSingleAnswer gets a single answer to a single question.
func FetchSingleAnswer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	qid := vars["id"]
	aid := vars["aid"]

	js, err := model.FetchSingleAnswer(qid, aid)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		// handle the error
		if err.Error() == "Answer not found" {
			w.WriteHeader(http.StatusNotFound)
			w.Write(js)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(js)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
