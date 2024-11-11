package sqlite

import (
	"strconv"

	"github.com/louisbranch/edulab"
	"github.com/pkg/errors"
)

func (db *DB) CreateQuestion(q *edulab.Question) error {
	query := `INSERT INTO questions (assessment_id, prompt, type)
	VALUES (?, ?, ?)`

	res, err := db.Exec(query, q.AssessmentID, q.Prompt, q.Type)
	if err != nil {
		return errors.Wrap(err, "could not create question")
	}

	id, err := res.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "retrieve last experiment id")
	}

	q.ID = strconv.FormatInt(id, 10)

	return nil
}

func (db *DB) CreateQuestionChoice(qc *edulab.QuestionChoice) error {
	query := `INSERT INTO question_choices (question_id, text, is_correct)
	VALUES (?, ?, ?)`

	_, err := db.Exec(query, qc.QuestionID, qc.Text, qc.IsCorrect)

	return errors.Wrap(err, "could not create question choice")
}

func (db *DB) FindQuestion(assessmentID string, pid string) (edulab.Question, error) {
	question := edulab.Question{
		AssessmentID: assessmentID,
	}

	query := `SELECT id, prompt, type
	FROM questions
	WHERE assessment_id AND id = ?`

	err := db.QueryRow(query, pid).Scan(&question.ID, &question.Prompt, &question.Type)
	if err != nil {
		return question, errors.Wrap(err, "could not find question")
	}
	return question, nil
}

func (db *DB) FindQuestions(assessmentID string) ([]edulab.Question, error) {

	query := `SELECT id, prompt, type
	FROM questions
	WHERE assessment_id = ?
	ORDER BY created_at ASC`

	rows, err := db.Query(query, assessmentID)
	if err != nil {
		return nil, errors.Wrap(err, "could not find questions")
	}

	defer rows.Close()

	var questions []edulab.Question
	for rows.Next() {
		q := edulab.Question{
			AssessmentID: assessmentID,
		}
		err = rows.Scan(&q.ID, &q.Prompt, &q.Type)
		if err != nil {
			return nil, errors.Wrap(err, "could not find questions")
		}

		questions = append(questions, q)
	}

	return questions, nil

}

func (db *DB) FindQuestionChoices(assessmentID string) ([]edulab.QuestionChoice, error) {

	query := `SELECT qc.id, qc.question_id, qc.text, qc.is_correct
	FROM question_choices AS qc
	JOIN questions AS q ON qc.question_id = q.id
	WHERE q.assessment_id = ?`

	rows, err := db.Query(query, assessmentID)
	if err != nil {
		return nil, errors.Wrap(err, "could not find question choices")
	}

	defer rows.Close()

	var choices []edulab.QuestionChoice
	for rows.Next() {
		c := edulab.QuestionChoice{}
		err = rows.Scan(&c.ID, &c.QuestionID, &c.Text, &c.IsCorrect)
		if err != nil {
			return nil, errors.Wrap(err, "could not find question choices")
		}

		choices = append(choices, c)
	}

	return choices, nil
}
