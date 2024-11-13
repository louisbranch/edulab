package postgres

import (
	"strconv"

	"github.com/pkg/errors"

	"github.com/louisbranch/edulab"
)

func (db *DB) CreateCohort(c *edulab.Cohort) error {
	query := `INSERT INTO cohorts (experiment_id, public_id, name, description)
		VALUES ($1, $2, $3, $4) RETURNING id`

	var id int64
	err := db.QueryRow(query, c.ExperimentID, c.PublicID, c.Name, c.Description).Scan(&id)
	if err != nil {
		return errors.Wrap(err, "could not create cohort")
	}

	c.ID = strconv.FormatInt(id, 10)

	return nil
}

func (db *DB) UpdateCohort(experimentID string, c edulab.Cohort) error {
	query := `UPDATE cohorts
		SET name = $1, description = $2
		WHERE experiment_id = $3 AND public_id = $4`

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
		WHERE experiment_id = $1 AND public_id = $2`

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
		WHERE experiment_id = $1
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
