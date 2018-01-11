package model

import (
	"encoding/json"
	"errors"
)

// FetchAllQuestions takes no parameters and responds with a JSON of all the questions in the db and an error.
func FetchAllQuestions() (js []byte, e error) {
	var questions []questionModel
	var _questions []transformedQuestion

	db.Find(&questions)

	if len(questions) <= 0 {
		err := errors.New("Not found")
		return []byte(""), err
	}

	for _, item := range questions {
		answeredBy := ""
		_questions = append(_questions, transformedQuestion{ID: item.ID, Question: item.Question, Answered: item.Answered, AnsweredBy: answeredBy})
	}

	js, err := json.Marshal(_questions)

	return js, err
}

// CreateQuestion takes in the request body and responds with JSON success message and an error.
func CreateQuestion(b []byte) (js []byte, e error) {
	// declare the data type that the data will be put into
	var question questionModel

	// Unmarshal the JSON formatted data b into the questionModel struct
	err := json.Unmarshal(b, &question)

	// Handle error if any
	if err != nil {
		return []byte("Something went wrong"), err
	}

	// Save the question to the DB
	db.Save(&question)

	// Return a success message (maybe edit later to return the question?)
	return []byte("{\"message\": \"Question successfully added\"}"), nil
}
