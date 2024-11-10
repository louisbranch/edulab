package sqlite

import (
	"strconv"

	"github.com/pkg/errors"

	"github.com/louisbranch/edulab"
)

func (db *DB) CreateExperiment(e *edulab.Experiment) error {
	q := `INSERT into experiments (public_id, name, description, created_at) values (?, ?, ?, ?);`

	res, err := db.Exec(q, e.PublicID, e.Name, e.Description, e.CreatedAt)
	if err != nil {
		return errors.Wrap(err, "create experiment")
	}

	id, err := res.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "retrieve last experiment id")
	}

	e.ID = strconv.FormatInt(id, 10)

	return nil
}

func (db *DB) FindExperiments() ([]edulab.Experiment, error) {
	var experiments []edulab.Experiment

	query := `SELECT id, public_id, name, description, created_at FROM experiments
    ORDER BY created_at DESC LIMIT 10
    `

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
	q := "SELECT id, name, description, created_at FROM experiments where public_id = ?"

	e := edulab.Experiment{
		PublicID: pid,
	}

	err := db.QueryRow(q, pid).Scan(&e.ID, &e.Name, &e.Description, &e.CreatedAt)

	if err != nil {
		return e, errors.Wrap(err, "find experiment")
	}

	return e, nil
}
