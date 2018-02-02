package model

import "fmt"

// FetchQuestionAnswers returns all answers for a given question
func FetchQuestionAnswers(id string) ([]TransformedAnswer, error) {
	var answers []AnswerModel
	var _answers []TransformedAnswer

	db.Find(&answers, "question_id = ?", id)

	if len(answers) <= 0 {
		return nil, ErrorNotFound
	}

	for _, item := range answers {
		_answers = append(_answers, TransformedAnswer{ID: item.ID, QuestionID: item.QuestionID, Answer: item.Answer, UserID: item.UserID, FName: item.FName, Cohort: item.Cohort})
	}

	return _answers, nil
}

// CreateAnswer takes the data from the controller and puts it into the database.
func CreateAnswer(qid int, ai AnswerModelInput) (TransformedAnswer, error) {
	var q QuestionModel

	a := AnswerModel{QuestionID: qid, Answer: ai.Answer, UserID: ai.UserID, FName: ai.FName, Cohort: ai.Cohort}

	// find the associated question and mark it answered
	db.Model(&q).Where("id = ?", qid).Update("answered", true)

	// save the answer in the DB
	db.Save(&a)

	_a := TransformedAnswer{QuestionID: a.QuestionID, ID: a.ID, Answer: a.Answer, UserID: a.UserID, FName: a.FName, Cohort: a.Cohort}

	return _a, nil
}

// FetchSingleAnswer takes the data from the controller and fetches a single answer matching url params.
func FetchSingleAnswer(qid string, aid string) (TransformedAnswer, error) {

	// Declare the data types to be used
	var a AnswerModel
	var _a TransformedAnswer

	// Conduct the database query
	db.Where("id = ? AND question_id = ?", aid, qid).First(&a)

	// Handle not found case
	if a.ID == 0 {
		return TransformedAnswer{}, ErrorNotFound
	}

	// Prepare answer to be sent back to controller
	_a = TransformedAnswer{ID: a.ID, QuestionID: a.QuestionID, Answer: a.Answer, UserID: a.UserID, FName: a.FName, Cohort: a.Cohort}
	return _a, nil
}

// UpdateAnswer updates the values of the answer in the db.
func UpdateAnswer(qid string, uid int, aid string, a AnswerModelInput) (TransformedAnswer, error) {
	// Declare the data types to be used
	var dbA AnswerModel
	var _a TransformedAnswer

	// Fetch the answer
	db.First(&dbA, "id = ?", aid)

	// Get the associated user (based on bearer token)
	var u UserModel
	db.First(&u, "id = ?", uid)

	// See if user is allowed to edit this question (either owner or admin)
	if a.UserID != uid && u.Admin == false {
		fmt.Println("this is the error")
		return TransformedAnswer{}, ErrorForbidden
	}

	// Update the answer in the database
	db.Model(&dbA).Update("answer", a.Answer)

	// Format answer for response
	_a = TransformedAnswer{ID: dbA.ID, QuestionID: dbA.QuestionID, Answer: dbA.Answer, FName: dbA.FName, Cohort: dbA.Cohort}

	return _a, nil
}
