package postgres

import (
	"strconv"

	"github.com/pkg/errors"

	"github.com/louisbranch/edulab"
)

func (db *DB) CreateAssessment(a *edulab.Assessment) error {
	query := `INSERT INTO assessments (experiment_id, public_id, description, type)
		VALUES ($1, $2, $3, $4) RETURNING id`

	var id int64
	err := db.QueryRow(query, a.ExperimentID, a.PublicID, a.Description, a.Type).Scan(&id)
	if err != nil {
		return errors.Wrap(err, "cannot create assessment")
	}

	a.ID = strconv.FormatInt(id, 10)

	return nil
}

func (db *DB) FindAssessment(parentID string, pid string) (edulab.Assessment, error) {
	q := `SELECT id, description, type
		FROM assessments WHERE experiment_id = $1 AND public_id = $2`

	e := edulab.Assessment{
		ExperimentID: parentID,
		PublicID:     pid,
	}

	err := db.QueryRow(q, parentID, pid).Scan(&e.ID, &e.Description, &e.Type)
	if err != nil {
		return e, errors.Wrap(err, "cannot find assessment")
	}

	return e, nil
}

func (db *DB) FindAssessments(experimentID string) ([]edulab.Assessment, error) {
	rows, err := db.Query(`
		SELECT a.id, a.experiment_id, a.public_id, a.description, a.type,
		COUNT(questions.id) AS q
		FROM assessments AS a
		LEFT JOIN questions ON questions.assessment_id = a.id
		WHERE a.experiment_id = $1
		GROUP BY a.id
		ORDER BY a.created_at ASC
	`, experimentID)
	if err != nil {
		return nil, errors.Wrap(err, "cannot query assessments")
	}
	defer rows.Close()

	var assessments []edulab.Assessment
	for rows.Next() {
		var a edulab.Assessment
		err = rows.Scan(&a.ID, &a.ExperimentID, &a.PublicID, &a.Description, &a.Type, &a.QuestionsCount)
		if err != nil {
			return nil, errors.Wrap(err, "cannot find assessments")
		}

		assessments = append(assessments, a)
	}

	return assessments, nil
}
