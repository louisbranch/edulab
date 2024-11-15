package postgres

import (
	"strconv"

	"github.com/pkg/errors"

	"github.com/louisbranch/edulab"
)

func (db *DB) CreateParticipant(p *edulab.Participant) error {
	q := `INSERT INTO participants (public_id, experiment_id, cohort_id, access_token)
		VALUES ($1, $2, $3, $4) RETURNING id`

	var id int64
	err := db.QueryRow(q, p.PublicID, p.ExperimentID, p.CohortID, p.AccessToken).Scan(&id)
	if err != nil {
		return errors.Wrap(err, "create participant")
	}

	p.ID = strconv.FormatInt(id, 10)

	return nil
}

func (db *DB) FindParticipant(experimentID string, accessToken string) (edulab.Participant, error) {
	var p edulab.Participant

	query := `SELECT id, public_id, experiment_id, cohort_id, access_token
		FROM participants WHERE experiment_id = $1 AND access_token = $2`

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
		FROM participants WHERE experiment_id = $1
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
