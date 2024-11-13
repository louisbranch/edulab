package edulab

import (
	"encoding/json"
	"time"
)

type Experiment struct {
	ID                string
	PublicID          string
	Name              string
	Description       string
	CreatedAt         time.Time
	ParticipantsCount int
}

type AssessmentType string

const (
	AssessmentTypePre AssessmentType = "pre"
	AssessmentTypePos AssessmentType = "post"
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
	InputSingle   InputType = "single"
	InputMultiple InputType = "multiple"
	InputText     InputType = "text"
)

type Question struct {
	ID           string
	AssessmentID string
	Text         string
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
	Text         string
	Type         InputType
}

type DemographicOption struct {
	ID            string
	DemographicID string
	Text          string
}

type Participant struct {
	ID           string
	PublicID     string
	ExperimentID string
	CohortID     string
	AccessToken  string
}

type Participation struct {
	ExperimentID  string
	AssessmentID  string
	ParticipantID string
	Answers       json.RawMessage `json:"answers"`
	Demographics  json.RawMessage `json:"demographics"`
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

	CreateDemographic(*Demographic) error
	FindDemographics(experimentID string) ([]Demographic, error)

	CreateDemographicOption(*DemographicOption) error
	FindDemographicOptions(experimentID string) ([]DemographicOption, error)

	CreateParticipant(*Participant) error
	FindParticipant(experimentID string, accessToken string) (Participant, error)
	FindParticipants(experimentID string) ([]Participant, error)

	CreateParticipation(*Participation) error
	UpdateParticipation(Participation) error
	FindParticipation(experimentID string, assessmentID string, participantID string) (Participation, error)
	FindParticipationsByParticipant(experimentID string, participantID string) ([]Participation, error)
	FindParticipationsByAssessment(experimentID string, assessmentID string) ([]Participation, error)
}
