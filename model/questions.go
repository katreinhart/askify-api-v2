package model

import (
	"encoding/json"
	"errors"
	"strconv"
)

// FetchAllQuestions takes no parameters and responds with a JSON of all the questions in the db and an error.
func FetchAllQuestions() (js []byte, e error) {

	// Declare data structures to hold data
	var questions []questionModel
	var _questions []transformedQuestion

	// Retrieve all questions from database
	db.Find(&questions)

	// If there are no questions, send back a 404
	if len(questions) <= 0 {
		err := errors.New("Not found")
		return []byte("{\"message\": \"Questions not found\"}"), err
	}

	// Transform each question into format to be sent back
	for _, item := range questions {
		uid, _ := strconv.Atoi(item.UserID)
		_questions = append(_questions, transformedQuestion{ID: item.ID, Question: item.Question, Answered: item.Answered, UserID: uid, FName: item.FName, Cohort: item.Cohort})
	}

	// Marshal data into json and return
	js, err := json.Marshal(_questions)
	return js, err
}

// CreateQuestion takes in the request body and responds with JSON success message and an error.
func CreateQuestion(b []byte, uid string) (js []byte, e error) {

	// declare the data type that the data will be put into
	var question questionModel

	// Unmarshal the JSON formatted data b into the questionModel struct
	err := json.Unmarshal(b, &question)
	question.UserID = uid

	// Handle error if any
	if err != nil {
		return []byte("{\"message\": \"Something went wrong.\"}"), err
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
	uid, _ := strconv.Atoi(question.UserID)
	_question := transformedQuestion{ID: question.ID, Question: question.Question, Answered: question.Answered, UserID: uid, FName: question.FName, Cohort: question.Cohort}

	// Marshal question into JS and return
	js, err := json.Marshal(_question)
	return js, err
}

// UpdateQuestion takes in an ID and a byte array, updates the appropriate database row, and returns a JSON formatted response and an error
func UpdateQuestion(id string, uid string, b []byte) ([]byte, error) {

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

	// Get the associated user (based on bearer token)
	var user userModel
	db.First(&user, "id = ?", uid)

	// See if user is allowed to edit this question (either owner or admin)
	if question.UserID != uid && user.Admin == false {
		err := errors.New("Unauthorized")
		return []byte("{\"message\": \"Not allowed\"}"), err
	}

	// Unmarshal the JSON from the request body into the updatedQuestion format
	err := json.Unmarshal(b, &updatedQuestion)

	// Handle JSON marshalling error
	if err != nil {
		err := errors.New("Update error")
		return []byte("{\"message\": \"Error updating the question\"}"), err
	}

	// Update the question and its answered status
	db.Model(&question).Update("question", updatedQuestion.Question)
	db.Model(&question).Update("answered", updatedQuestion.Answered)

	// Format question for response
	_question = transformedQuestion{ID: updatedQuestion.ID, Question: updatedQuestion.Question, Answered: updatedQuestion.Answered, FName: updatedQuestion.FName, Cohort: updatedQuestion.Cohort}

	// Marshal into JSON and return
	js, err := json.Marshal(_question)
	return js, err
}

// FetchQueue will return all unanswered questions in the proper order
func FetchQueue() ([]byte, error) {

	// Data structures to hold the questions
	var questions []questionModel
	var _questions []transformedQuestion

	// Database call for unanswered questions, ordered by created_at timestamps
	db.Order("created_at asc").Find(&questions, "answered = ?", false)

	// Handle no questions returned from DB
	if len(questions) <= 0 {
		err := errors.New("Not found")
		return []byte("{\"message\": \"No unanswered questions found.\"}"), err
	}

	// Transform data into return format
	for _, item := range questions {
		uid, _ := strconv.Atoi(item.UserID)
		_questions = append(_questions, transformedQuestion{ID: item.ID, Question: item.Question, Answered: item.Answered, UserID: uid, FName: item.FName, Cohort: item.Cohort})
	}

	// Marshal into JSON and return
	js, err := json.Marshal(_questions)
	return js, err
}

// FetchArchive will return a nested object with all answered questions and their answers.
func FetchArchive() ([]byte, error) {
	var questions []questionModel
	var _questions []archiveQuestion

	// Get all answered questions from the database.
	db.Limit(50).Order("created_at desc").Find(&questions, "answered = ?", true)

	// No questions found? return not found error
	if len(questions) <= 0 {
		err := errors.New("Not found")
		return []byte("{\"message\": \"No questions found.\"}"), err
	}

	// go through each question and put its answers into an []archiveAnswer
	for _, q := range questions {
		var answers []answerModel
		var _answers []archiveAnswer

		// Get answers from the database
		db.Find(&answers, "question_id = ?", q.ID)

		// go through each answer, find the user, and put the answer in the []archiveAnswer slice
		for _, a := range answers {
			var user userModel
			db.First(&user, "id = ?", a.UserID)
			_answers = append(_answers, archiveAnswer{ID: a.ID, QuestionID: int(q.ID), Answer: a.Answer, UserID: a.UserID, FName: user.FName, Cohort: user.Cohort})
		}

		// find the user who posted the question
		var u userModel
		db.First(&u, "id = ?", q.UserID)

		// put the question and the []archiveAnswers into an []archiveQuestion slice
		uid, _ := strconv.Atoi(q.UserID)
		_questions = append(_questions, archiveQuestion{ID: q.ID, Question: q.Question, UserID: uid, FName: u.FName, Answers: _answers, Cohort: u.Cohort})
	}

	// marshal the data into JSON & return
	js, err := json.Marshal(_questions)
	return js, err
}

// FetchUserQuestions returns all questions a user has asked
func FetchUserQuestions(uid string) ([]byte, error) {
	var questions []questionModel
	var _questions []transformedQuestion

	// Database call for user's questions
	db.Order("created_at asc").Find(&questions, "user_id = ?", uid)

	// Handle no questions found
	if len(questions) <= 0 {
		err := errors.New("Not found")
		return []byte("{\"message\": \"No questions found for user.\"}"), err
	}

	// transform questions into format for sending back to FE
	for _, item := range questions {
		uid, _ := strconv.Atoi(item.UserID)
		_questions = append(_questions, transformedQuestion{ID: item.ID, Question: item.Question, Answered: item.Answered, UserID: uid, FName: item.FName, Cohort: item.Cohort})
	}

	// Marshal into JSON and return
	js, err := json.Marshal(_questions)
	return js, err
}
