package models

import (
	"math/rand"

	"github.com/louisbranch/edulab"
)

type Model struct {
	rand *rand.Rand
	db   edulab.Database
}

func New(db edulab.Database, rand *rand.Rand) Model {
	return Model{
		db:   db,
		rand: rand,
	}
}

func (m Model) CreateExperiment(e *edulab.Experiment) error {
	e.PublicID = m.NewPublicID([]int{3, 3})

	err := m.db.CreateExperiment(e)
	if err != nil {
		return err
	}

	err = m.db.CreateAssessment(&edulab.Assessment{
		ExperimentID: e.ID,
		PublicID:     m.NewPublicID([]int{3}),
		Type:         edulab.PreAssessment,
	})
	if err != nil {
		return err
	}

	err = m.db.CreateAssessment(&edulab.Assessment{
		ExperimentID: e.ID,
		PublicID:     m.NewPublicID([]int{3}),
		Type:         edulab.PostAssessment,
	})
	if err != nil {
		return err
	}

	return nil
}

func (m Model) CreateCohort(cohort *edulab.Cohort) error {
	cohort.PublicID = m.NewPublicID([]int{3})

	return m.db.CreateCohort(cohort)
}
