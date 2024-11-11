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
	ID             string
	ExperimentID   string
	PublicID       string
	Name           string
	Description    string
	IsPre          bool
	QuestionsCount int
}

<<<<<<< Updated upstream
type AssessmentQuestion struct {
=======
type QuestionType string

const (
	SingleChoice   QuestionType = "single_choice"
	MultipleChoice QuestionType = "multiple_choice"
	FreeForm       QuestionType = "free_form"
)

type Question struct {
>>>>>>> Stashed changes
	ID           string
	AssessmentID string
	Text         string
	Type         string
}

<<<<<<< Updated upstream
type AssessmentChoice struct {
	ID            string
	AssessmentID  string
	AssessmentQID string
	Text          string
	IsCorrect     bool
=======
type QuestionChoice struct {
	ID         string
	QuestionID string
	Text       string
	IsCorrect  bool
>>>>>>> Stashed changes
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
	FindExperiment(publicID string) (Experiment, error)

	CreateAssessment(*Assessment) error
<<<<<<< Updated upstream
	FindAssessments(string) ([]Assessment, error)
=======
	FindAssessment(experimentID string, publicID string) (Assessment, error)
	FindAssessments(experimentID string) ([]Assessment, error)

	CreateQuestion(*Question) error
	FindQuestion(assessmentID string, publicID string) (Question, error)
	FindQuestions(assessmentID string) ([]Question, error)

	CreateQuestionChoice(*QuestionChoice) error
>>>>>>> Stashed changes
}
