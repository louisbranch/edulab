package wizard

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"math/rand"
	"sort"

	"github.com/pkg/errors"

	"github.com/louisbranch/edulab"
)

type BootstrapConfig struct {
	Participants      int                `yaml:"participants"` // Total number of participants to create
	AssessmentConfigs []AssessmentConfig `yaml:"assessments"`  // Configuration for each assessment
	DemographicConfig DemographicConfig  `yaml:"demographics"` // Configuration for demographics
}

type DemographicConfig struct {
	Probabilities      []float64 `yaml:"probabilities"`       // Probabilities for each option in order
	OutlierProbability float64   `yaml:"outlier_probability"` // Probability of overriding with a random outlier
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

	// Fetch demographics for the experiment
	demographics, err := db.FindDemographics(experiment.ID)
	if err != nil {
		return errors.Wrap(err, "could not find demographics for experiment")
	}
	demographicOptions := make(map[string][]edulab.DemographicOption)

	options, err := db.FindDemographicOptions(experiment.ID)
	if err != nil {
		return errors.Wrapf(err, "could not find options for experiment %s", experiment.ID)
	}
	for _, option := range options {
		id := option.DemographicID
		demographicOptions[id] = append(demographicOptions[id], option)
	}

	// Assume a single demographics config to apply to all fields
	demographicConfig := config.DemographicConfig

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

		// Generate demographics for each participant
		demographicResponses := make(map[string]string)
		for _, demographic := range demographics {
			options := demographicOptions[demographic.ID]
			selectedOption := weightedRandomChoice(options, demographicConfig.Probabilities,
				demographicConfig.OutlierProbability)
			demographicResponses[demographic.ID] = selectedOption.ID
		}

		// Serialize demographics to JSON
		demographicsJSON, err := json.Marshal(demographicResponses)
		if err != nil {
			return errors.Wrap(err, "could not marshal demographic responses to JSON")
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

			if j == 0 {
				participation.Demographics = demographicsJSON
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
		selected := selectChoices(question, choices, correctProbability, biasFactor)
		sort.Strings(selected)
		answers[question.ID] = selected
	}

	return json.Marshal(answers)
}

func selectChoices(question edulab.Question, choices []edulab.QuestionChoice,
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

	seededRand := rand.New(rand.NewSource(hashStringToSeed(question.ID)))

	switch question.Type {
	case "single":
		// Single choice: Use correctProbability to pick correct vs. incorrect
		if rand.Float64() < correctProbability {
			// Pick the correct choice
			if len(correctChoices) > 0 {
				selectedChoices = append(selectedChoices, correctChoices[0])
			}
		} else {
			// Pick a "biased" incorrect choice using seeded random for consistency
			if rand.Float64() < biasFactor && len(incorrectChoices) > 0 {
				// Seeded random generator to consistently select the same preferred incorrect choice per question
				selectedChoices = append(selectedChoices, incorrectChoices[seededRand.Intn(len(incorrectChoices))])
			} else if len(incorrectChoices) > 0 {
				// Use any random incorrect choice with the global random generator
				selectedChoices = append(selectedChoices, incorrectChoices[rand.Intn(len(incorrectChoices))])
			}
		}

	case "multiple":
		// Multiple choice: Use correctProbability to pick all correct or a mix
		if rand.Float64() < correctProbability {
			// Pick all correct choices
			selectedChoices = append(selectedChoices, correctChoices...)
		} else {
			// Pick a random subset of correct and/or incorrect choices
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

func weightedRandomChoice(options []edulab.DemographicOption,
	probabilities []float64, outlierProbability float64) edulab.DemographicOption {
	if len(options) == 0 {
		return edulab.DemographicOption{} // No options available
	}

	// Ensure probabilities match the number of options, default to 0 for extra options
	if len(probabilities) < len(options) {
		diff := len(options) - len(probabilities)
		probabilities = append(probabilities, make([]float64, diff)...)
	}

	// Normalize probabilities to sum to 1
	total := 0.0
	for _, p := range probabilities {
		total += p
	}
	if total > 0 {
		for i := range probabilities {
			probabilities[i] /= total
		}
	}

	// Outlier check: override with a random outlier if probability is met
	if rand.Float64() < outlierProbability {
		return options[rand.Intn(len(options))]
	}

	// Weighted random selection
	r := rand.Float64()
	cumulative := 0.0
	for i, option := range options {
		cumulative += probabilities[i]
		if r < cumulative {
			return option
		}
	}

	// Fallback (should not be reached if probabilities are normalized)
	return options[0]
}

// Helper function to convert a string (question ID) to a consistent seed value
func hashStringToSeed(s string) int64 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int64(h.Sum32())
}
