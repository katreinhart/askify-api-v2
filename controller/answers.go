package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/katreinhart/askify-api-v2/model"
)

// FetchQuestionAnswers fetches all the answers to a given question
func FetchQuestionAnswers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	qid := vars["id"]

	var answers []model.TransformedAnswer

	answers, err := model.FetchQuestionAnswers(qid)

	if err != nil {
		handleErrorAndRespond(nil, err, w)
		return
	}

	js, err := json.Marshal(answers)

	handleErrorAndRespond(js, err, w)
}

// CreateAnswer posts an answer to a question
func CreateAnswer(w http.ResponseWriter, r *http.Request) {
	// get the data from the request body
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	var a model.AnswerModelInput
	var _a model.TransformedAnswer

	err := json.Unmarshal(b, &a)
	if err != nil {
		handleErrorAndRespond(nil, model.ErrorBadRequest, w)
		return
	}

	vars := mux.Vars(r)
	qid, _ := strconv.Atoi(vars["id"])

	// parse UID from bearer token
	uid, err := GetUIDFromBearerToken(r)
	if err != nil {
		handleErrorAndRespond(nil, err, w)
		return
	}

	a.UserID, _ = strconv.Atoi(uid)

	_a, err = model.CreateAnswer(qid, a)

	if err != nil {
		handleErrorAndRespond(nil, model.ErrorInternalServer, w)
		return
	}

	js, err := json.Marshal(_a)
	handleErrorAndRespond(js, err, w)
}

// FetchSingleAnswer gets a single answer to a single question.
func FetchSingleAnswer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	qid := vars["id"]
	aid := vars["aid"]

	var a model.TransformedAnswer
	a, err := model.FetchSingleAnswer(qid, aid)

	if err != nil {
		handleErrorAndRespond(nil, err, w)
		return
	}

	js, err := json.Marshal(a)

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
		handleErrorAndRespond(nil, err, w)
		return
	}

	var a model.AnswerModelInput

	err = json.Unmarshal(b, &a)
	if err != nil {
		handleErrorAndRespond(nil, model.ErrorBadRequest, w)
	}

	userID, _ := strconv.Atoi(uid)

	_a, err := model.UpdateAnswer(qid, userID, aid, a)

	if err != nil {
		handleErrorAndRespond(nil, err, w)
		return
	}

	js, err := json.Marshal(_a)

	handleErrorAndRespond(js, err, w)
}
