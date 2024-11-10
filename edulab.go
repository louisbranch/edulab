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

type QuestionType string

const (
	SingleChoice   QuestionType = "single_choice"
	MultipleChoice QuestionType = "multiple_choice"
	FreeForm       QuestionType = "free_form"
)

type AssessmentQuestion struct {
	ID           string
	AssessmentID string
	Prompt       string
	Type         QuestionType
}

type AssessmentChoice struct {
	ID            string
	AssessmentID  string
	AssessmentQID string
	Value         string
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
	FindAssessment(string, string) (Assessment, error)
	FindAssessments(string) ([]Assessment, error)
}
