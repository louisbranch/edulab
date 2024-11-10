package edulab

import "time"

type Experiment struct {
	ID          string
	PublicID    string
	Name        string
	Description string
	CreatedAt   time.Time
}

type Assessment struct {
	ID           string
	ExperimentID string
	Name         string
	Description  string
	IsPre        bool
	CreatedAt    time.Time
}

type Database interface {
	CreateExperiment(*Experiment) error
	UpdateExperiment(Experiment) error
	FindExperiments() ([]Experiment, error)
	FindExperiment(string) (Experiment, error)

	CreateAssessment(*Assessment) error
	FindAssessments(string) ([]Assessment, error)
}
