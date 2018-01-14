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
		_questions = append(_questions, transformedQuestion{ID: item.ID, Question: item.Question, Answered: item.Answered, UserID: item.UserID})
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

// FetchSingleQuestion takes in an ID and returns a JSON formatted question and an error
func FetchSingleQuestion(id string) ([]byte, error) {
	// Declare the data type that will hold the data
	var question questionModel
	// Fetch the question from the DB
	db.First(&question, id)

	// This is the error handling for GORM
	if question.ID == 0 {
		err := errors.New("Not found")
		return []byte("{\"message\": \"Question not found\"}"), err
	}

	// Transform the question into the response type
	_question := transformedQuestion{ID: question.ID, Question: question.Question, Answered: question.Answered, UserID: question.UserID}

	js, err := json.Marshal(_question)

	return js, err
}

// UpdateQuestion takes in an ID and a byte array, updates the appropriate database row, and returns a JSON formatted response and an error
func UpdateQuestion(id string, b []byte) ([]byte, error) {
	// Declare the data types to be used
	var question, updatedQuestion questionModel
	var _question transformedQuestion
	// Fetch the question in question
	db.First(&question, id)

	// Handle not found error
	if question.ID == 0 {
		err := errors.New("Not found")
		return []byte("{\"message\": \"Question not found\"}"), err
	}

	// Unmarshal the JSON from the request body into the updatedQuestion format
	err := json.Unmarshal(b, &updatedQuestion)
	if err != nil {
		// hanlde the error
		err := errors.New("Update error")
		return []byte("{\"message\": \"Error updating the question\"}"), err
	}

	db.Model(&question).Update("question", updatedQuestion.Question)
	db.Model(&question).Update("answered", updatedQuestion.Answered)

	_question = transformedQuestion{ID: updatedQuestion.ID, Question: updatedQuestion.Question, Answered: updatedQuestion.Answered}
	js, err := json.Marshal(_question)

	return js, err
}

// FetchQueue will return all unanswered questions in the proper order
func FetchQueue() ([]byte, error) {
	var questions []questionModel
	var _questions []transformedQuestion

	db.Order("created_at asc").Find(&questions, "answered = ?", false)

	if len(questions) <= 0 {
		err := errors.New("Not found")
		return []byte(""), err
	}

	for _, item := range questions {
		_questions = append(_questions, transformedQuestion{ID: item.ID, Question: item.Question, Answered: item.Answered, UserID: item.UserID})
	}

	js, err := json.Marshal(_questions)

	return js, err
}

// FetchArchive will return a nested object with all answered questions and their answers.
func FetchArchive() ([]byte, error) {
	var questions []questionModel
	var _questions []archiveQuestion

	db.Order("created_at asc").Find(&questions, "answered = ?", true)

	if len(questions) <= 0 {
		err := errors.New("Not found")
		return []byte(""), err
	}

	for _, q := range questions {
		var answers []answerModel
		var _answers []archiveAnswer
		db.Find(&answers, "question_id = ?", q.ID)
		for _, a := range answers {
			var user userModel
			db.First(&user, "id = ?", a.UserID)
			_answers = append(_answers, archiveAnswer{ID: a.ID, QuestionID: int(q.ID), Answer: a.Answer, UserID: a.UserID, UserName: user.Fname})
		}
		var u userModel
		db.First(&u, "id = ?", q.UserID)
		_questions = append(_questions, archiveQuestion{ID: q.ID, Question: q.Question, UserID: q.UserID, UserName: u.Fname, Answers: _answers})
	}

	js, err := json.Marshal(_questions)

	return js, err
}

// FetchUserQuestions returns all questions a user has asked
func FetchUserQuestions(uid string) ([]byte, error) {
	var questions []questionModel
	var _questions []transformedQuestion

	db.Order("created_at asc").Find(&questions, "user_id = ?", uid)

	if len(questions) <= 0 {
		err := errors.New("Not found")
		return []byte(""), err
	}

	for _, item := range questions {
		_questions = append(_questions, transformedQuestion{ID: item.ID, Question: item.Question, Answered: item.Answered, UserID: item.UserID})
	}

	js, err := json.Marshal(_questions)

	return js, err
}
