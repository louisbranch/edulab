package result

import (
	"encoding/json"
	"errors"

	"github.com/louisbranch/edulab"
)

// QuestionScore calculates the score for each participation for a given question
// and returns a map from cohortID to score.
func (r *Result) QuestionScore(questionID string) (map[string][]float64, error) {
	question, exists := r.questions[questionID]
	if !exists {
		return nil, errors.New("question not found")
	}

	// Map to store scores by cohort ID
	scores := make(map[string][]float64)
	for participantID, participations := range r.participation {
		for _, participation := range participations {
			// Deserialize answers for this participation
			var answers map[string][]string
			if err := json.Unmarshal(participation.Answers, &answers); err != nil {
				return nil, err
			}

			// Get participant's answer for the question, if available
			answerIDs, answered := answers[questionID]
			if !answered {
				continue // Skip if participant didn't answer this question
			}

			participant := r.participants[participantID]
			cohortID := participant.CohortID

			// Score the answer based on question type
			var score float64
			switch question.Type {
			case edulab.InputSingle:
				score = r.scoreSingleAnswer(questionID, answerIDs)
			case edulab.InputMultiple:
				score = r.scoreMultipleAnswer(questionID, answerIDs)
			}

			// Append score to cohort's list of scores
			if _, exists := scores[cohortID]; !exists {
				scores[cohortID] = []float64{}
			}

			scores[cohortID] = append(scores[cohortID], score)
		}
	}
	return scores, nil
}

// scoreSingleAnswer scores a single-answer question as 0 or 1.
func (r *Result) scoreSingleAnswer(questionID string, answerIDs []string) float64 {
	if len(answerIDs) != 1 {
		return 0.0 // Invalid answer for single-input type
	}
	correctChoices := r.getCorrectChoices(questionID)
	for _, choice := range correctChoices {
		if choice.ID == answerIDs[0] {
			return 1.0
		}
	}
	return 0.0
}

// scoreMultipleAnswer scores a multiple-answer question from 0 to 1.
func (r *Result) scoreMultipleAnswer(questionID string, answerIDs []string) float64 {
	correctChoices := r.getCorrectChoices(questionID)
	correctSet := make(map[string]bool)
	for _, choice := range correctChoices {
		correctSet[choice.ID] = true
	}

	// Calculate ratio of correct answers given
	total := len(correctSet)
	if total == 0 {
		return 0.0
	}

	correct := 0
	wrong := 0
	for _, answerID := range answerIDs {
		if correctSet[answerID] {
			correct++
		} else {
			wrong++
		}
	}

	// Calculate the points per correct answer
	points := 1.0 / float64(total)

	// Calculate score based on correct and incorrect answers
	score := float64(correct)*points - float64(wrong)*points

	// Ensure score is not negative
	if score < 0 {
		return 0
	}
	return score
}

// getCorrectChoices retrieves the correct choices for a question.
func (r *Result) getCorrectChoices(questionID string) []edulab.QuestionChoice {
	var correctChoices []edulab.QuestionChoice
	for _, choice := range r.choices[questionID] {
		if choice.IsCorrect {
			correctChoices = append(correctChoices, choice)
		}
	}
	return correctChoices
}
