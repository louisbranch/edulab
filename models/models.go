package models

import (
	"math/rand"

	"golang.org/x/text/message"

	"github.com/louisbranch/edulab"
)

type Model struct {
	rand    *rand.Rand
	printer *message.Printer
	db      edulab.Database
}

func New(db edulab.Database, rand *rand.Rand, printer *message.Printer) Model {
	return Model{
		db:      db,
		rand:    rand,
		printer: printer,
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
		Name:         m.printer.Sprintf("Pre-Assessment"),
		IsPre:        true,
	})
	if err != nil {
		return err
	}

	err = m.db.CreateAssessment(&edulab.Assessment{
		ExperimentID: e.ID,
		PublicID:     m.NewPublicID([]int{3}),
		Name:         m.printer.Sprintf("Post-Assessment"),
		IsPre:        false,
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
