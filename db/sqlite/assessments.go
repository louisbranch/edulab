package sqlite

import (
	"github.com/louisbranch/edulab"
	"github.com/pkg/errors"
)

func (db *DB) CreateAssessment(a *edulab.Assessment) error {
	_, err := db.Exec(`
		INSERT INTO assessments (experiment_id, public_id, name, description, is_pre)
		VALUES (?, ?, ?, ?, ?)
	`, a.ExperimentID, a.PublicID, a.Name, a.Description, a.IsPre)

	return err
}

func (db *DB) FindAssessment(parentID string, pid string) (edulab.Assessment, error) {
	q := `SELECT id, name, description, is_pre
	FROM assessments where experiment_id = ? AND public_id = ?`

	e := edulab.Assessment{
		ExperimentID: parentID,
		PublicID:     pid,
	}

	err := db.QueryRow(q, parentID, pid).Scan(&e.ID, &e.Name, &e.Description, &e.IsPre)

	if err != nil {
		return e, errors.Wrap(err, "find experiment")
	}

	return e, nil
}

func (db *DB) FindAssessments(experimentID string) ([]edulab.Assessment, error) {
	rows, err := db.Query(`
		SELECT id, experiment_id, public_id, name, description, is_pre
		FROM assessments
		WHERE experiment_id = ?
		ORDER BY created_at ASC
	`, experimentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assessments []edulab.Assessment
	for rows.Next() {
		var a edulab.Assessment
		err = rows.Scan(&a.ID, &a.ExperimentID, &a.PublicID, &a.Name, &a.Description, &a.IsPre)
		if err != nil {
			return nil, err
		}

		assessments = append(assessments, a)
	}

	return assessments, nil
}
