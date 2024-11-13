package sqlite

import (
	"database/sql"

	"github.com/pkg/errors"

	"github.com/louisbranch/edulab"
)

func (db *DB) CreateParticipation(p *edulab.Participation) error {
	q := `INSERT into participations (experiment_id,
	assessment_id, participant_id, answers, demographics)
	values (?, ?, ?, ?, ?);`

	_, err := db.Exec(q, p.ExperimentID, p.AssessmentID, p.ParticipantID,
		p.Answers, p.Demographics)
	if err != nil {
		return errors.Wrap(err, "create participation")
	}

	return nil
}

func (db *DB) UpdateParticipation(p edulab.Participation) error {
	q := `UPDATE participations SET answers = ?, demographics = ?
	WHERE experiment_id = ? AND assessment_id = ? AND participant_id = ?`

	_, err := db.Exec(q, p.Answers, p.Demographics, p.ExperimentID, p.AssessmentID, p.ParticipantID)
	if err != nil {
		return errors.Wrap(err, "update participation")
	}

	return nil
}

func (db *DB) FindParticipation(experimentID, assessmentID, participantID string) (edulab.Participation, error) {
	var p edulab.Participation

	query := `SELECT experiment_id, assessment_id, participant_id, answers, demographics
	FROM participations WHERE experiment_id = ? AND assessment_id = ? AND participant_id = ?`

	var answers, demographics sql.NullString

	err := db.QueryRow(query, experimentID, assessmentID, participantID).
		Scan(&p.ExperimentID, &p.AssessmentID, &p.ParticipantID, &answers, &demographics)
	if err != nil {
		return p, errors.Wrap(err, "query participation")
	}

	if answers.Valid {
		p.Answers = []byte(answers.String)
	}
	if demographics.Valid {
		p.Demographics = []byte(demographics.String)
	}

	return p, nil
}

func (db *DB) FindParticipationsByParticipant(experimentID, participantID string) ([]edulab.Participation, error) {
	var participations []edulab.Participation

	query := `SELECT assessment_id, answers, demographics
	FROM participations WHERE experiment_id = ? AND participant_id = ?`

	rows, err := db.Query(query, experimentID, participantID)
	if err != nil {
		return nil, errors.Wrap(err, "query participations")
	}
	defer rows.Close()

	for rows.Next() {
		p := edulab.Participation{
			ExperimentID:  experimentID,
			ParticipantID: participantID,
		}

		var answers, demographics sql.NullString
		err = rows.Scan(&p.AssessmentID, &answers, &demographics)
		if err != nil {
			return nil, errors.Wrap(err, "scan participation")
		}

		if answers.Valid {
			p.Answers = []byte(answers.String)
		}
		if demographics.Valid {
			p.Demographics = []byte(demographics.String)
		}

		participations = append(participations, p)
	}

	return participations, nil
}

func (db *DB) FindParticipationsByAssessment(experimentID, assessmentID string) ([]edulab.Participation, error) {

	var participations []edulab.Participation

	query := `SELECT participant_id, answers, demographics
	FROM participations WHERE experiment_id = ? AND assessment_id = ?`

	rows, err := db.Query(query, experimentID, assessmentID)
	if err != nil {
		return nil, errors.Wrap(err, "query participations")
	}
	defer rows.Close()

	for rows.Next() {
		p := edulab.Participation{
			ExperimentID: experimentID,
			AssessmentID: assessmentID,
		}

		var answers, demographics sql.NullString

		err = rows.Scan(&p.ParticipantID, &answers, &demographics)
		if err != nil {
			return nil, errors.Wrap(err, "scan participation")
		}

		if answers.Valid {
			p.Answers = []byte(answers.String)
		}
		if demographics.Valid {
			p.Demographics = []byte(demographics.String)
		}

		participations = append(participations, p)
	}

	return participations, nil
}
