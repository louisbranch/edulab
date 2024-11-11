package sqlite

import "github.com/louisbranch/edulab"

func (db *DB) CreateAssessment(a *edulab.Assessment) error {
	query := `INSERT INTO assessments (experiment_id, public_id, name, description, is_pre)
		VALUES (?, ?, ?, ?, ?)`

	_, err := db.Exec(query, a.ExperimentID, a.PublicID, a.Name, a.Description, a.IsPre)

	return errors.Wrap(err, "cannot create assessment")
}

<<<<<<< Updated upstream
=======
func (db *DB) FindAssessment(parentID string, pid string) (edulab.Assessment, error) {
	q := `SELECT id, name, description, is_pre
	FROM assessments where experiment_id = ? AND public_id = ?`

	e := edulab.Assessment{
		ExperimentID: parentID,
		PublicID:     pid,
	}

	err := db.QueryRow(q, parentID, pid).Scan(&e.ID, &e.Name, &e.Description, &e.IsPre)

	if err != nil {
		return e, errors.Wrap(err, "cannot find assessment")
	}

	return e, nil
}

>>>>>>> Stashed changes
func (db *DB) FindAssessments(experimentID string) ([]edulab.Assessment, error) {
	rows, err := db.Query(`
		SELECT assessments.id, experiment_id, assessments.public_id, name, description, is_pre,
		COUNT(questions.id) AS questions
		FROM assessments
		LEFT JOIN questions ON questions.assessment_id = assessments.id
		WHERE experiment_id = ?
		GROUP BY assessments.id
		ORDER BY assessments.created_at ASC
	`, experimentID)
	if err != nil {
		return nil, errors.Wrap(err, "cannot query assessments")
	}
	defer rows.Close()

	var assessments []edulab.Assessment
	for rows.Next() {
		var a edulab.Assessment
		err = rows.Scan(&a.ID, &a.ExperimentID, &a.PublicID, &a.Name, &a.Description,
			&a.IsPre, &a.QuestionsCount)
		if err != nil {
			return nil, errors.Wrap(err, "cannot find assessments")
		}

		assessments = append(assessments, a)
	}

	return assessments, nil
}
