package mock

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/louisbranch/edulab"
)

type DB struct {
	experiments        []edulab.Experiment
	assessments        []edulab.Assessment
	questions          []edulab.Question
	questionChoices    []edulab.QuestionChoice
	cohorts            []edulab.Cohort
	demographics       []edulab.Demographic
	demographicOptions []edulab.DemographicOption
	participants       []edulab.Participant
	participations     []edulab.Participation
}

func NewDB() *DB {
	db := &DB{}
	db.LoadFixtures(filepath.Join("..", "mock", "fixtures"))
	return db
}

// LoadFixtures loads YAML fixtures into the mock database
func (db *DB) LoadFixtures(fixturesDir string) error {
	// Load experiments
	if err := loadYAML(filepath.Join(fixturesDir, "experiments.yaml"), &db.experiments); err != nil {
		return err
	}
	// Load assessments
	if err := loadYAML(filepath.Join(fixturesDir, "assessments.yaml"), &db.assessments); err != nil {
		return err
	}
	// Load questions
	if err := loadYAML(filepath.Join(fixturesDir, "questions.yaml"), &db.questions); err != nil {
		return err
	}
	// Load question choices
	if err := loadYAML(filepath.Join(fixturesDir, "question_choices.yaml"), &db.questionChoices); err != nil {
		return err
	}
	// Load cohorts
	if err := loadYAML(filepath.Join(fixturesDir, "cohorts.yaml"), &db.cohorts); err != nil {
		return err
	}
	// Load demographics
	if err := loadYAML(filepath.Join(fixturesDir, "demographics.yaml"), &db.demographics); err != nil {
		return err
	}
	// Load demographic options
	if err := loadYAML(filepath.Join(fixturesDir, "demographic_options.yaml"), &db.demographicOptions); err != nil {
		return err
	}
	// Load participants
	if err := loadYAML(filepath.Join(fixturesDir, "participants.yaml"), &db.participants); err != nil {
		return err
	}
	// Load participations
	if err := loadYAML(filepath.Join(fixturesDir, "participations.yaml"), &db.participations); err != nil {
		return err
	}
	return nil
}

// loadYAML is a helper function to load a YAML file into a given struct
func loadYAML(filePath string, out interface{}) error {
	data, err := os.ReadFile(filePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return yaml.Unmarshal(data, out)
}

// CreateExperiment creates a new experiment
func (db *DB) CreateExperiment(e *edulab.Experiment) error {
	db.experiments = append(db.experiments, *e)
	return nil
}

// FindExperiments fetches all experiments
func (db *DB) FindExperiments() ([]edulab.Experiment, error) {
	return db.experiments, nil
}

// FindExperiment fetches an experiment by public ID
func (db *DB) FindExperiment(publicID string) (edulab.Experiment, error) {
	for _, e := range db.experiments {
		if e.PublicID == publicID {
			return e, nil
		}
	}
	return edulab.Experiment{}, errors.New("no experiment found")
}

// UpdateExperiment updates an existing experiment
func (db *DB) UpdateExperiment(e edulab.Experiment) error {
	for i, ex := range db.experiments {
		if ex.PublicID == e.PublicID {
			db.experiments[i] = e
			return nil
		}
	}
	return errors.New("no experiment found")
}

// DeleteExperiment deletes an existing experiment
func (db *DB) DeleteExperiment(publicID string) error {
	for i, e := range db.experiments {
		if e.PublicID == publicID {
			db.experiments = append(db.experiments[:i], db.experiments[i+1:]...)
			return nil
		}
	}
	return errors.New("no experiment found")
}

// CreateAssessment creates a new assessment
func (db *DB) CreateAssessment(a *edulab.Assessment) error {
	db.assessments = append(db.assessments, *a)
	return nil
}

// FindAssessment fetches an assessment by public ID
func (db *DB) FindAssessment(experimentID, publicID string) (edulab.Assessment, error) {
	for _, a := range db.assessments {
		if a.ExperimentID == experimentID && a.PublicID == publicID {
			return a, nil
		}
	}
	return edulab.Assessment{}, errors.New("no assessment found")
}

// FindAssessments fetches assessments by experiment ID
func (db *DB) FindAssessments(experimentID string) ([]edulab.Assessment, error) {
	var result []edulab.Assessment
	for _, a := range db.assessments {
		if a.ExperimentID == experimentID {
			result = append(result, a)
		}
	}
	return result, nil
}

// CreateQuestion creates a new question
func (db *DB) CreateQuestion(q *edulab.Question) error {
	db.questions = append(db.questions, *q)
	return nil
}

// FindQuestion fetches a question by public ID
func (db *DB) FindQuestion(assessmentID, id string) (edulab.Question, error) {
	for _, q := range db.questions {
		if q.AssessmentID == assessmentID && q.ID == id {
			return q, nil
		}
	}

	return edulab.Question{}, errors.New("no question found")
}

// FindQuestions fetches questions by assessment ID
func (db *DB) FindQuestions(assessmentID string) ([]edulab.Question, error) {
	var result []edulab.Question
	for _, q := range db.questions {
		if q.AssessmentID == assessmentID {
			result = append(result, q)
		}
	}
	return result, nil
}

// CreateQuestionChoice creates a new question choice
func (db *DB) CreateQuestionChoice(qc *edulab.QuestionChoice) error {
	db.questionChoices = append(db.questionChoices, *qc)
	return nil
}

// FindQuestionChoices fetches question choices by question ID
func (db *DB) FindQuestionChoices(assessmentID string) ([]edulab.QuestionChoice, error) {
	questions := make(map[string]string)
	for _, q := range db.questions {
		questions[q.ID] = q.AssessmentID
	}

	var result []edulab.QuestionChoice
	for _, c := range db.questionChoices {
		if questions[c.QuestionID] == assessmentID {
			result = append(result, c)
		}
	}
	return result, nil
}

// CreateCohort creates a new cohort
func (db *DB) CreateCohort(c *edulab.Cohort) error {
	db.cohorts = append(db.cohorts, *c)
	return nil
}

// UpdateCohort updates an existing cohort
func (db *DB) UpdateCohort(experimentID string, c edulab.Cohort) error {
	for i, co := range db.cohorts {
		if co.ExperimentID == experimentID && co.ID == c.ID {
			db.cohorts[i] = c
			return nil
		}
	}
	return errors.New("no cohort found")
}

// FindCohort fetches a cohort by public ID
func (db *DB) FindCohort(experimentID, publicID string) (edulab.Cohort, error) {
	for _, c := range db.cohorts {
		if c.ExperimentID == experimentID && c.PublicID == publicID {
			return c, nil
		}
	}
	return edulab.Cohort{}, errors.New("no cohort found")
}

// FindCohorts fetches cohorts by experiment ID
func (db *DB) FindCohorts(experimentID string) ([]edulab.Cohort, error) {
	var result []edulab.Cohort
	for _, c := range db.cohorts {
		if c.ExperimentID == experimentID {
			result = append(result, c)
		}
	}
	return result, nil
}

// CreateDemographic creates a new demographic
func (db *DB) CreateDemographic(d *edulab.Demographic) error {
	db.demographics = append(db.demographics, *d)
	return nil
}

// CreateDemographicOption creates a new demographic option
func (db *DB) CreateDemographicOption(d *edulab.DemographicOption) error {
	db.demographicOptions = append(db.demographicOptions, *d)
	return nil
}

// FindDemographics fetches demographics by experiment ID
func (db *DB) FindDemographics(experimentID string) ([]edulab.Demographic, error) {
	var result []edulab.Demographic
	for _, d := range db.demographics {
		if d.ExperimentID == experimentID {
			result = append(result, d)
		}
	}
	return result, nil
}

// FindDemographicOptions fetches demographic options by demographic ID
func (db *DB) FindDemographicOptions(demographicID string) ([]edulab.DemographicOption, error) {
	var result []edulab.DemographicOption
	for _, d := range db.demographicOptions {
		if d.DemographicID == demographicID {
			result = append(result, d)
		}
	}
	return result, nil
}

// CreateParticipant creates a new participant
func (db *DB) CreateParticipant(p *edulab.Participant) error {
	db.participants = append(db.participants, *p)
	return nil
}

// FindParticipant fetches a participant by public ID
func (db *DB) FindParticipant(experimentID, accessToken string) (edulab.Participant, error) {
	for _, p := range db.participants {
		if p.ExperimentID == experimentID && p.AccessToken == accessToken {
			return p, nil
		}
	}
	return edulab.Participant{}, errors.New("no participant found")
}

// FindParticipants fetches participants by experiment ID
func (db *DB) FindParticipants(experimentID string) ([]edulab.Participant, error) {
	var result []edulab.Participant
	for _, p := range db.participants {
		if p.ExperimentID == experimentID {
			result = append(result, p)
		}
	}
	return result, nil
}

// CreateParticipation creates a new participation
func (db *DB) CreateParticipation(p *edulab.Participation) error {
	db.participations = append(db.participations, *p)
	return nil
}

// UpdateParticipation updates an existing participation
func (db *DB) UpdateParticipation(p edulab.Participation) error {
	for i, pa := range db.participations {
		if pa.ExperimentID == p.ExperimentID && pa.CohortID == p.CohortID && pa.ParticipantID == p.ParticipantID {
			db.participations[i] = p
			return nil
		}
	}
	return errors.New("no participation found")
}

// FindParticipation fetches a participation by public ID
func (db *DB) FindParticipation(experimentID, assessmentID, participantID string) (edulab.Participation, error) {
	for _, p := range db.participations {
		if p.ExperimentID == experimentID && p.AssessmentID == assessmentID && p.ParticipantID == participantID {
			return p, nil
		}
	}
	return edulab.Participation{}, errors.New("no participation found")
}

// FindParticipations fetches participations by experiment ID
func (db *DB) FindParticipations(experimentID string) ([]edulab.Participation, error) {
	var result []edulab.Participation
	for _, p := range db.participations {
		if p.ExperimentID == experimentID {
			result = append(result, p)
		}
	}
	return result, nil
}

// FindParticipationsByAssessment fetches participations by assessment ID
func (db *DB) FindParticipationsByAssessment(experimentID, assessmentID string) ([]edulab.Participation, error) {
	var result []edulab.Participation
	for _, p := range db.participations {
		if p.ExperimentID == experimentID && p.AssessmentID == assessmentID {
			result = append(result, p)
		}
	}
	return result, nil
}

// FindParticipationsByParticipant fetches participations by participant ID
func (db *DB) FindParticipationsByParticipant(experimentID, participantID string) ([]edulab.Participation, error) {
	var result []edulab.Participation
	for _, p := range db.participations {
		if p.ExperimentID == experimentID && p.ParticipantID == participantID {
			result = append(result, p)
		}
	}
	return result, nil
}
