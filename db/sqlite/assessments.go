package sqlite

import "github.com/louisbranch/edulab"

func (db *DB) CreateAssessment(a *edulab.Assessment) error {
	_, err := db.Exec(`
		INSERT INTO assessments (experiment_id, public_id, name, description, is_pre)
		VALUES (?, ?, ?, ?, ?)
	`, a.ExperimentID, a.PublicID, a.Name, a.Description, a.IsPre)

	return err
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
