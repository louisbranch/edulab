package edulab

import "time"

type Experiment struct {
	ID          string
	PublicID    string
	Name        string
	Description string
	CreatedAt   time.Time
}

type AssessmentType string

const (
	PreAssessment  AssessmentType = "pre_assessment"
	PostAssessment AssessmentType = "post_assessment"
)

type Assessment struct {
	ID             string
	ExperimentID   string
	PublicID       string
	Description    string
	Type           AssessmentType
	QuestionsCount int
}

type InputType string

const (
	SingleChoice   InputType = "single_choice"
	MultipleChoice InputType = "multiple_choice"
	Text           InputType = "text"
)

type Question struct {
	ID           string
	AssessmentID string
	Prompt       string
	Type         InputType
}

type QuestionChoice struct {
	ID         string
	QuestionID string
	Text       string
	IsCorrect  bool
}

type Cohort struct {
	ID           string
	ExperimentID string
	PublicID     string
	Name         string
	Description  string
}

type Demographic struct {
	ID           string
	ExperimentID string
	Translatable
	Type InputType
}

type DemographicOption struct {
	ID            string
	DemographicID string
	Translatable
}

type Translatable struct {
	I18n string
	Text string
}

type Database interface {
	CreateExperiment(*Experiment) error
	UpdateExperiment(Experiment) error
	FindExperiments() ([]Experiment, error)
	FindExperiment(publicID string) (Experiment, error)

	CreateAssessment(*Assessment) error
	FindAssessment(experimentID string, publicID string) (Assessment, error)
	FindAssessments(experimentID string) ([]Assessment, error)

	CreateQuestion(*Question) error
	FindQuestion(assessmentID string, publicID string) (Question, error)
	FindQuestions(assessmentID string) ([]Question, error)

	CreateQuestionChoice(*QuestionChoice) error
	FindQuestionChoices(assessmentID string) ([]QuestionChoice, error)

	CreateCohort(*Cohort) error
	UpdateCohort(experimentID string, c Cohort) error
	FindCohort(experimentID string, publicID string) (Cohort, error)
	FindCohorts(experimentID string) ([]Cohort, error)
}
