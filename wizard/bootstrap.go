package wizard

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sort"

	"github.com/pkg/errors"

	"github.com/louisbranch/edulab"
)

type BootstrapConfig struct {
	Participants      int                `yaml:"participants"` // Total number of participants to create
	AssessmentConfigs []AssessmentConfig `yaml:"assessments"`  // Configuration for each assessment in the experiment
}

type AssessmentConfig struct {
	CorrectProbabilities []float64 `yaml:"correct_probabilities"` // Probability of correct answers per cohort for this assessment
	BiasFactor           float64   `yaml:"bias_factor"`           // Likelihood of selecting a common incorrect answer when wrong
}

func bootstrapParticipants(db edulab.Database, config BootstrapConfig,
	experiment edulab.Experiment) error {

	if config.Participants <= 0 {
		log.Printf("[INFO] No participants to bootstrap for experiment %q.\n", experiment.Name)
		return nil
	}

	participants, err := db.FindParticipants(experiment.ID)
	if err != nil {
		return errors.Wrap(err, "could not find participants")
	}

	// Fetch the cohorts for this experiment
	cohorts, err := db.FindCohorts(experiment.ID)
	if err != nil {
		return errors.Wrap(err, "could not find cohorts for experiment")
	}
	numCohorts := len(cohorts)

	// Now simulate responses for each assessment
	assessments, err := db.FindAssessments(experiment.ID)
	if err != nil {
		return errors.Wrap(err, "could not find assessments")
	}

	// Validate configuration
	if len(config.AssessmentConfigs) != len(assessments) {
		return errors.New("incorrect number of assessment configurations")
	}

	for i, assessmentConfig := range config.AssessmentConfigs {
		if len(assessmentConfig.CorrectProbabilities) != numCohorts {
			return errors.Errorf("incorrect number of correct probabilities for assessment %d", i+1)
		}
	}

	for i := len(participants); i < config.Participants; i++ {
		cohort := cohorts[i%numCohorts] // Distribute participants evenly across cohorts
		participant := edulab.Participant{
			PublicID:     fmt.Sprintf("%s-P%d", experiment.PublicID, i+1),
			ExperimentID: experiment.ID,
			CohortID:     cohort.ID,
			AccessToken:  fmt.Sprintf("token-%s-%d", experiment.PublicID, i+1),
		}

		// Create the participant in the database
		if err := db.CreateParticipant(&participant); err != nil {
			return errors.Wrap(err, "could not create participant")
		}

		for j, assessment := range assessments {
			// Retrieve configuration for this assessment
			assessmentConfig := config.AssessmentConfigs[j]
			cohortProb := assessmentConfig.CorrectProbabilities[i%numCohorts]
			correctProbability := cohortProb

			// Generate answers with bias applied
			answers, err := generateAnswers(db, assessment, correctProbability, assessmentConfig.BiasFactor)
			if err != nil {
				return errors.Wrapf(err, "could not generate answers for participant %s", participant.PublicID)
			}

			// Save the participation record
			participation := edulab.Participation{
				ExperimentID:  experiment.ID,
				AssessmentID:  assessment.ID,
				ParticipantID: participant.ID,
				Answers:       answers,
			}

			if err := db.CreateParticipation(&participation); err != nil {
				return errors.Wrap(err, "could not create participation record")
			}
		}
	}

	needed := config.Participants - len(participants)
	log.Printf("[INFO] Successfully bootstrapped %d participants for experiment %q.\n",
		needed, experiment.Name)
	return nil
}

func generateAnswers(db edulab.Database, assessment edulab.Assessment,
	correctProbability float64, biasFactor float64) ([]byte, error) {

	questions, err := db.FindQuestions(assessment.ID)
	if err != nil {
		return nil, errors.Wrap(err, "could not find questions for assessment")
	}
	allChoices, err := db.FindQuestionChoices(assessment.ID)
	if err != nil {
		return nil, errors.Wrap(err, "could not find question choices")
	}

	answers := make(map[string][]string)

	for _, question := range questions {

		var choices []edulab.QuestionChoice
		for _, choice := range allChoices {
			if choice.QuestionID == question.ID {
				choices = append(choices, choice)
			}
		}

		// Add selected choices to answers map
		selected := selectChoices(question.Type, choices, correctProbability, biasFactor)
		sort.Strings(selected)
		answers[question.ID] = selected
	}

	return json.Marshal(answers)
}

func selectChoices(inputType edulab.InputType, choices []edulab.QuestionChoice,
	correctProbability float64, biasFactor float64) []string {
	var correctChoices, incorrectChoices, selectedChoices []string

	// Separate choices into correct and incorrect lists
	for _, choice := range choices {
		if choice.IsCorrect {
			correctChoices = append(correctChoices, choice.ID)
		} else {
			incorrectChoices = append(incorrectChoices, choice.ID)
		}
	}

	switch inputType {
	case edulab.InputSingle:
		// Single choice: Use correctProbability to pick correct vs. incorrect
		if rand.Float64() < correctProbability {
			// Pick the correct choice
			if len(correctChoices) > 0 {
				selectedChoices = append(selectedChoices, correctChoices[0])
			}
		} else {
			// Pick an incorrect choice based on biasFactor
			if rand.Float64() < biasFactor && len(incorrectChoices) > 0 {
				selectedChoices = append(selectedChoices, incorrectChoices[rand.Intn(len(incorrectChoices))])
			} else if len(incorrectChoices) > 0 {
				selectedChoices = append(selectedChoices, incorrectChoices[rand.Intn(len(incorrectChoices))])
			}
		}

	case edulab.InputMultiple:
		// Multiple choice: Use correctProbability to pick all correct or a mix
		if rand.Float64() < correctProbability {
			// Pick all correct choices
			selectedChoices = append(selectedChoices, correctChoices...)
		} else {
			// Pick a random subset of correct and/or incorrect choices based on biasFactor
			numChoices := rand.Intn(len(choices) + 1) // Pick any number of choices from 0 to all
			allChoices := append(correctChoices, incorrectChoices...)
			rand.Shuffle(len(allChoices), func(i, j int) {
				allChoices[i], allChoices[j] = allChoices[j], allChoices[i]
			})
			selectedChoices = append(selectedChoices, allChoices[:numChoices]...)
		}
	}

	return selectedChoices
}
