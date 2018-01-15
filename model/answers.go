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
		_answers = append(_answers, transformedAnswer{ID: item.ID, QuestionID: item.QuestionID, Answer: item.Answer, UserID: item.UserID, UserFName: item.UserFName, Cohort: item.Cohort})
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

	_answer = answerModel{QuestionID: qid, Answer: answer.Answer, UserID: answer.UserID, UserFName: answer.UserFName, Cohort: answer.Cohort}

	if err != nil {
		return []byte("{\"message\": \"Error saving to database\"}"), err
	}

	// find the associated question and mark it answered
	db.Model(&question).Update("answered", true)

	// save the answer in the DB
	db.Save(&_answer)

	return []byte("{\"message\": \"Answer created successfully\"}"), nil
}

// FetchSingleAnswer takes the data from the controller and fetches a single answer matching url params.
func FetchSingleAnswer(qid string, aid string) ([]byte, error) {
	var answer answerModel
	var _answer transformedAnswer

	db.Find(&answer, "id = ?", aid)

	if answer.ID == 0 {
		return []byte("{\"message\": \"Answer not found\"}"), errors.New("Not found")
	}

	_answer = transformedAnswer{ID: answer.ID, QuestionID: answer.QuestionID, Answer: answer.Answer, UserID: answer.UserID, UserFName: answer.UserFName, Cohort: answer.Cohort}

	js, err := json.Marshal(_answer)

	return js, err
}
