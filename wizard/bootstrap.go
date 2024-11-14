package wizard

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"

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

	answers := make(map[string]interface{})

	for _, question := range questions {
		choices, err := db.FindQuestionChoices(question.ID)
		if err != nil {
			return nil, errors.Wrap(err, "could not find question choices")
		}

		// Determine if the participant answers correctly based on the probability
		isCorrect := rand.Float64() < correctProbability
		var selectedChoices []string

		if isCorrect {
			// Select only the correct choices
			for _, choice := range choices {
				if choice.IsCorrect {
					selectedChoices = append(selectedChoices, choice.ID)
				}
			}
		} else {
			// Apply bias to select common incorrect choices
			commonIncorrectChoices := getCommonIncorrectChoices(choices)
			if rand.Float64() < biasFactor && len(commonIncorrectChoices) > 0 {
				// Pick from common incorrect choices if bias applies
				selectedChoices = []string{commonIncorrectChoices[rand.Intn(len(commonIncorrectChoices))]}
			} else {
				// Otherwise, pick randomly from other incorrect choices
				selectedChoices = randomIncorrectChoices(choices)
			}
		}

		// Add selected choices to answers map
		answers[question.ID] = selectedChoices
	}

	return json.Marshal(answers)
}

// Helper function to identify common incorrect choices
func getCommonIncorrectChoices(choices []edulab.QuestionChoice) []string {
	var commonIncorrectChoices []string
	for _, choice := range choices {
		if !choice.IsCorrect && isCommonMisconception(choice) {
			commonIncorrectChoices = append(commonIncorrectChoices, choice.ID)
		}
	}
	return commonIncorrectChoices
}

// Function to tag certain choices as common misconceptions
func isCommonMisconception(choice edulab.QuestionChoice) bool {
	// Example: Tag choices with known misconceptions (customize as needed)
	return choice.Text == "The **distance** between the Earth and the Sun changes throughout the year."
}

// Helper to pick a random incorrect choice (excluding common ones)
func randomIncorrectChoices(choices []edulab.QuestionChoice) []string {
	var incorrectChoices []string
	for _, choice := range choices {
		if !choice.IsCorrect {
			incorrectChoices = append(incorrectChoices, choice.ID)
		}
	}
	// Pick one random incorrect choice
	if len(incorrectChoices) > 0 {
		return []string{incorrectChoices[rand.Intn(len(incorrectChoices))]}
	}
	return nil
}
