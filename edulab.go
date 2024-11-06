package edulab

import "time"

type Experiment struct {
	ID           string
	Name         string
	Participants int
	CreatedAt    time.Time
}

type Database interface {
	CreateExperiment(*Experiment) error
	FindExperiments() ([]Experiment, error)
	FindExperiment(string) (Experiment, error)
}
