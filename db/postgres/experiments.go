package postgres

import (
	"errors"

	"github.com/louisbranch/edulab"
)

func (db *DB) CreateExperiment(e *edulab.Experiment) error {
	return errors.New("not implemented")
}

func (db *DB) UpdateExperiment(e edulab.Experiment) error {
	return errors.New("not implemented")
}

func (db *DB) FindExperiments() ([]edulab.Experiment, error) {
	return nil, errors.New("not implemented")
}

func (db *DB) FindExperiment(name string) (edulab.Experiment, error) {
	e := edulab.Experiment{}

	return e, errors.New("not implemented")
}
