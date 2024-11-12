package sqlite

import (
	"strconv"

	"github.com/pkg/errors"

	"github.com/louisbranch/edulab"
)

func (db *DB) CreateDemographic(d *edulab.Demographic) error {
	query := `INSERT INTO demographics (experiment_id, text, type)
	VALUES (?, ?, ?)`

	res, err := db.Exec(query, d.ExperimentID, d.Text, d.Type)
	if err != nil {
		return errors.Wrap(err, "could not create demographic")
	}

	id, err := res.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "retrieve last demographic id")
	}

	d.ID = strconv.FormatInt(id, 10)

	return nil
}

func (db *DB) FindDemographics(experimentID string) ([]edulab.Demographic, error) {
	var demographics []edulab.Demographic

	query := `SELECT id, text, type
	FROM demographics
	WHERE experiment_id = ?`

	rows, err := db.Query(query, experimentID)
	if err != nil {
		return demographics, errors.Wrap(err, "could not find demographics")
	}
	defer rows.Close()

	for rows.Next() {
		var d edulab.Demographic
		err := rows.Scan(&d.ID, &d.Text, &d.Type)
		if err != nil {
			return demographics, errors.Wrap(err, "could not scan demographic")
		}
		d.ExperimentID = experimentID

		demographics = append(demographics, d)
	}

	return demographics, nil
}

func (db *DB) CreateDemographicOption(o *edulab.DemographicOption) error {
	query := `INSERT INTO demographic_options (demographic_id, text) VALUES (?, ?)`

	res, err := db.Exec(query, o.DemographicID, o.Text)
	if err != nil {
		return errors.Wrap(err, "could not create demographic option")
	}

	id, err := res.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "retrieve last demographic option id")
	}

	o.ID = strconv.FormatInt(id, 10)

	return nil
}

func (db *DB) FindDemographicOptions(experimentID string) ([]edulab.DemographicOption, error) {
	var options []edulab.DemographicOption

	query := `SELECT o.id, o.demographic_id, o.text
	FROM demographic_options AS o
	JOIN demographics AS d ON o.demographic_id = d.id
	WHERE d.experiment_id = ?`

	rows, err := db.Query(query, experimentID)
	if err != nil {
		return options, errors.Wrap(err, "could not find demographic options")
	}
	defer rows.Close()

	for rows.Next() {
		var o edulab.DemographicOption
		err := rows.Scan(&o.ID, &o.DemographicID, &o.Text)
		if err != nil {
			return options, errors.Wrap(err, "could not scan demographic option")
		}

		options = append(options, o)
	}

	return options, nil
}
