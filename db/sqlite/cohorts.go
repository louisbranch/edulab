package sqlite

import (
	"strconv"

	"github.com/pkg/errors"

	"github.com/louisbranch/edulab"
)

func (db *DB) CreateCohort(c *edulab.Cohort) error {
	query := `INSERT INTO cohorts (experiment_id, public_id, name, description)
	VALUES (?, ?, ?, ?)`

	res, err := db.Exec(query, c.ExperimentID, c.PublicID, c.Name, c.Description)
	if err != nil {
		return errors.Wrap(err, "could not create cohort")
	}

	id, err := res.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "retrieve last cohort id")
	}

	c.ID = strconv.FormatInt(id, 10)

	return nil
}

func (db *DB) UpdateCohort(experimentID string, c edulab.Cohort) error {
	query := `UPDATE cohorts
	SET name = ?, description = ?
	WHERE experiment_id = ? AND public_id = ?`

	_, err := db.Exec(query, c.Name, c.Description, experimentID, c.PublicID)
	if err != nil {
		return errors.Wrap(err, "could not update cohort")
	}
	return nil
}

func (db *DB) FindCohort(experimentID string, publicID string) (edulab.Cohort, error) {
	cohort := edulab.Cohort{
		ExperimentID: experimentID,
		PublicID:     publicID,
	}

	query := `SELECT id, name, description
	FROM cohorts
	WHERE experiment_id = ? AND public_id = ?`

	err := db.QueryRow(query, experimentID, publicID).Scan(&cohort.ID, &cohort.Name, &cohort.Description)
	if err != nil {
		return cohort, errors.Wrap(err, "could not find cohort")
	}
	return cohort, nil
}

func (db *DB) FindCohorts(experimentID string) ([]edulab.Cohort, error) {
	var cohorts []edulab.Cohort

	query := `SELECT id, public_id, name, description
	FROM cohorts
	WHERE experiment_id = ?
	ORDER BY created_at ASC`

	rows, err := db.Query(query, experimentID)
	if err != nil {
		return cohorts, errors.Wrap(err, "could not find cohorts")
	}
	defer rows.Close()

	for rows.Next() {
		cohort := edulab.Cohort{
			ExperimentID: experimentID,
		}
		err := rows.Scan(&cohort.ID, &cohort.PublicID, &cohort.Name, &cohort.Description)
		if err != nil {
			return cohorts, errors.Wrap(err, "could not scan cohort")
		}
		cohorts = append(cohorts, cohort)
	}

	return cohorts, nil
}
