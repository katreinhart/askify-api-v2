package model

import (
	"encoding/json"
	"errors"
)

// FetchQuestionAnswers returns all answers for a given question
func FetchQuestionAnswers(id string) ([]byte, error) {
	var answers []answerModel
	var _answers []transformedAnswer

	db.Find(&answers, "question_id = ?", id)

	if len(answers) <= 0 {
		return []byte("{\"message\": \"Answers not found\"}"), errors.New("Not found")
	}

	for _, item := range answers {
		_answers = append(_answers, transformedAnswer{ID: item.ID, QuestionID: item.QuestionID, Answer: item.Answer, UserID: item.UserID, FName: item.FName, Cohort: item.Cohort})
	}

	js, err := json.Marshal(_answers)

	return js, err
}

// CreateAnswer takes the data from the controller and puts it into the database.
func CreateAnswer(qid int, b []byte) ([]byte, error) {
	var answer answerModelInput
	var _answer answerModel

	var question questionModel

	err := json.Unmarshal(b, &answer)

	_answer = answerModel{QuestionID: qid, Answer: answer.Answer, UserID: answer.UserID, FName: answer.FName, Cohort: answer.Cohort}

	if err != nil {
		return []byte("{\"message\": \"Error saving to database\"}"), err
	}

	// find the associated question and mark it answered
	db.Model(&question).Where("id = ?", qid).Update("answered", true)

	// save the answer in the DB
	db.Save(&_answer)

	return []byte("{\"message\": \"Answer created successfully\"}"), nil
}

// FetchSingleAnswer takes the data from the controller and fetches a single answer matching url params.
func FetchSingleAnswer(qid string, aid string) ([]byte, error) {

	// Declare the data types to be used
	var answer answerModel
	var _answer transformedAnswer

	// Conduct the database query
	db.Where("id = ? AND question_id = ?", aid, qid).First(&answer)

	// Handle not found case
	if answer.ID == 0 {
		return []byte("{\"message\": \"Answer not found\"}"), errors.New("Not found")
	}

	// Prepare answer to be sent; marshal into JSON and return
	_answer = transformedAnswer{ID: answer.ID, QuestionID: answer.QuestionID, Answer: answer.Answer, UserID: answer.UserID, FName: answer.FName, Cohort: answer.Cohort}
	js, err := json.Marshal(_answer)
	return js, err
}

// UpdateAnswer updates the values of the answer in the db.
func UpdateAnswer(qid string, uid string, aid string, b []byte) ([]byte, error) {
	// Declare the data types to be used
	var answer, updatedAnswer answerModel
	var _answer transformedAnswer

	// Fetch the answer
	db.First(&answer, "id = ?", aid)

	// Get the associated user (based on bearer token)
	var user userModel
	db.First(&user, "id = ?", uid)

	// See if user is allowed to edit this question (either owner or admin)
	if string(answer.UserID) != uid && user.Admin == false {
		err := errors.New("Unauthorized")
		return []byte("{\"message\": \"Not allowed\"}"), err
	}

	// Unmarshal the JSON from the request body into the updatedAnswer format
	err := json.Unmarshal(b, &updatedAnswer)

	// Handle JSON marshalling error
	if err != nil {
		err := errors.New("Update error")
		return []byte("{\"message\": \"Error updating the answer\"}"), err
	}

	// Update the answer in the database
	db.Model(&answer).Update("answer", updatedAnswer.Answer)

	// Format answer for response
	_answer = transformedAnswer{ID: updatedAnswer.ID, Answer: updatedAnswer.Answer, FName: updatedAnswer.FName, Cohort: updatedAnswer.Cohort}

	// Marshal into JSON and return
	js, err := json.Marshal(_answer)
	return js, err
}
