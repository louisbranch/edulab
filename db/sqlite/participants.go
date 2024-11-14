package sqlite

import (
	"strconv"

	"github.com/pkg/errors"

	"github.com/louisbranch/edulab"
)

func (db *DB) CreateParticipant(p *edulab.Participant) error {
	q := `INSERT into participants (public_id, experiment_id, cohort_id, access_token)
	values (?, ?, ?, ?);`

	res, err := db.Exec(q, p.PublicID, p.ExperimentID, p.CohortID, p.AccessToken)
	if err != nil {
		return errors.Wrap(err, "create participant")
	}

	id, err := res.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "retrieve last participant id")
	}

	p.ID = strconv.FormatInt(id, 10)

	return nil
}

func (db *DB) FindParticipant(experimentID string, accessToken string) (edulab.Participant, error) {
	var p edulab.Participant

	query := `SELECT id, public_id, experiment_id, cohort_id, access_token
	FROM participants WHERE experiment_id = ? AND access_token = ?`

	err := db.QueryRow(query, experimentID, accessToken).
		Scan(&p.ID, &p.PublicID, &p.ExperimentID, &p.CohortID, &p.AccessToken)
	if err != nil {
		return p, errors.Wrap(err, "query participant")
	}

	return p, nil
}

func (db *DB) FindParticipants(experimentID string) ([]edulab.Participant, error) {
	var participants []edulab.Participant

	query := `SELECT id, public_id, experiment_id, cohort_id, access_token
	FROM participants WHERE experiment_id = ?
	ORDER BY id`

	rows, err := db.Query(query, experimentID)
	if err != nil {
		return nil, errors.Wrap(err, "query participants")
	}
	defer rows.Close()

	for rows.Next() {
		p := edulab.Participant{}
		err = rows.Scan(&p.ID, &p.PublicID, &p.ExperimentID, &p.CohortID, &p.AccessToken)
		if err != nil {
			return nil, errors.Wrap(err, "scan participants")
		}
		participants = append(participants, p)
	}

	return participants, nil
}
