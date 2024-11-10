package postgres

import (
	"errors"

	"github.com/louisbranch/edulab"
)

func (db *DB) CreateAssessment(a *edulab.Assessment) error {
	return errors.New("not implemented")
}

func (db *DB) FindAssessments(experimentID string) ([]edulab.Assessment, error) {
	return nil, errors.New("not implemented")
}
