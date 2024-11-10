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
	PublicID     string
	Name         string
	Description  string
	IsPre        bool
}

type AssessmentQuestion struct {
	ID           string
	AssessmentID string
	Text         string
	Type         string
}

type AssessmentChoice struct {
	ID            string
	AssessmentID  string
	AssessmentQID string
	Text          string
	IsCorrect     bool
}

type Cohort struct {
	ID           string
	ExperimentID string
	PublicID     string
	Name         string
	Description  string
}

type Database interface {
	CreateExperiment(*Experiment) error
	UpdateExperiment(Experiment) error
	FindExperiments() ([]Experiment, error)
	FindExperiment(string) (Experiment, error)

	CreateAssessment(*Assessment) error
	FindAssessments(string) ([]Assessment, error)
}
