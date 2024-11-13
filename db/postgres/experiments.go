package postgres

import (
	"strconv"

	"github.com/pkg/errors"

	"github.com/louisbranch/edulab"
)

func (db *DB) CreateExperiment(e *edulab.Experiment) error {
	q := `INSERT INTO experiments (public_id, name, description) VALUES ($1, $2, $3) RETURNING id`

	var id int64
	err := db.QueryRow(q, e.PublicID, e.Name, e.Description).Scan(&id)
	if err != nil {
		return errors.Wrap(err, "create experiment")
	}

	e.ID = strconv.FormatInt(id, 10)

	return nil
}

func (db *DB) UpdateExperiment(e edulab.Experiment) error {
	q := `UPDATE experiments SET name = $1, description = $2 WHERE public_id = $3`

	_, err := db.Exec(q, e.Name, e.Description, e.PublicID)
	if err != nil {
		return errors.Wrap(err, "update experiment")
	}

	return nil
}

func (db *DB) FindExperiments() ([]edulab.Experiment, error) {
	var experiments []edulab.Experiment

	query := `SELECT id, public_id, name, description, created_at 
		FROM experiments
		ORDER BY created_at DESC LIMIT 10`

	rows, err := db.Query(query)
	if err != nil {
		return nil, errors.Wrap(err, "query experiments")
	}
	defer rows.Close()

	for rows.Next() {
		e := edulab.Experiment{}
		err = rows.Scan(&e.ID, &e.PublicID, &e.Name, &e.Description, &e.CreatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "scan experiments")
		}
		experiments = append(experiments, e)
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, "find experiments")
	}
	return experiments, nil
}

func (db *DB) FindExperiment(pid string) (edulab.Experiment, error) {
	q := `SELECT id, name, description, created_at FROM experiments WHERE public_id = $1`

	e := edulab.Experiment{
		PublicID: pid,
	}

	err := db.QueryRow(q, pid).Scan(&e.ID, &e.Name, &e.Description, &e.CreatedAt)
	if err != nil {
		return e, errors.Wrap(err, "find experiment")
	}

	return e, nil
}
