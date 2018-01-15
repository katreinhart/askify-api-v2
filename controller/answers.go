package controller

import (
	"bytes"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/katreinhart/askify-api-v2/model"
)

// FetchQuestionAnswers fetches all the answers to a given question
func FetchQuestionAnswers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	qid := vars["id"]

	js, err := model.FetchQuestionAnswers(qid)

	handleErrorAndRespond(js, err, w)
}

// CreateAnswer posts an answer to a question
func CreateAnswer(w http.ResponseWriter, r *http.Request) {
	// get the data from the request body
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	vars := mux.Vars(r)
	qid, _ := strconv.Atoi(vars["id"])

	js, err := model.CreateAnswer(qid, b)

	handleErrorAndRespond(js, err, w)
}

// FetchSingleAnswer gets a single answer to a single question.
func FetchSingleAnswer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	qid := vars["id"]
	aid := vars["aid"]

	js, err := model.FetchSingleAnswer(qid, aid)

	handleErrorAndRespond(js, err, w)
}

// UpdateAnswer handles PUT requests to /questions/{id}/answers/{aid}
func UpdateAnswer(w http.ResponseWriter, r *http.Request) {
	// get the data from the request body
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	vars := mux.Vars(r)
	qid := vars["id"]
	aid := vars["aid"]

	// parse UID from bearer token
	uid, err := GetUIDFromBearerToken(r)
	if err != nil {
		handleErrorAndRespond([]byte("{\"message\": \"Error parsing bearer token.\"}"), err, w)
	}

	js, err := model.UpdateAnswer(qid, uid, aid, b)

	handleErrorAndRespond(js, err, w)
}
