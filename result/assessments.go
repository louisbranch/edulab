package result

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/louisbranch/edulab"
)

func CountChoicesByCohorts(db edulab.Database, experiment edulab.Experiment) ([][][]int, error) {

	participants, err := db.FindParticipants(experiment.ID)
	if err != nil {
		return nil, err
	}

	participantCohorts := make(map[string]string)
	for _, p := range participants {
		participantCohorts[p.ID] = p.CohortID
	}

	participations, err := db.FindParticipations(experiment.ID)
	if err != nil {
		return nil, err
	}

	cohorts, err := db.FindCohorts(experiment.ID)
	if err != nil {
		return nil, err
	}

	questionsMap := make(map[string]map[string]map[string]int)

	for _, p := range participations {

		if p.Answers == nil {
			continue
		}

		if p.CohortID == "" {
			p.CohortID = participantCohorts[p.ParticipantID]
		}

		var answers map[string][]string
		if err := json.Unmarshal(p.Answers, &answers); err != nil {
			return nil, errors.Wrapf(err, "failed to unmarshal answers for participation %s", p.Answers)
		}

		for questionID, choiceIDs := range answers {
			if _, ok := questionsMap[questionID]; !ok {
				questionsMap[questionID] = make(map[string]map[string]int)
			}

			if _, ok := questionsMap[questionID][p.CohortID]; !ok {
				questionsMap[questionID][p.CohortID] = make(map[string]int)
			}

			for _, choiceID := range choiceIDs {
				questionsMap[questionID][p.CohortID][choiceID]++
			}
		}
	}

	assessments, err := db.FindAssessments(experiment.ID)
	if err != nil {
		return nil, err
	}

	var total [][][]int

	for _, assessment := range assessments {

		questions, err := db.FindQuestions(assessment.ID)
		if err != nil {
			return nil, err
		}

		choices, err := db.FindQuestionChoices(assessment.ID)
		if err != nil {
			return nil, err
		}

		for _, question := range questions {

			counts := make([][]int, len(cohorts))

			for _, choice := range choices {
				if choice.QuestionID != question.ID {
					continue
				}

				if _, ok := questionsMap[choice.QuestionID]; !ok {
					continue
				}

				for i, cohort := range cohorts {
					count := questionsMap[choice.QuestionID][cohort.ID][choice.ID]
					counts[i] = append(counts[i], count)
				}

			}

			add := false
			for _, count := range counts {
				for _, c := range count {
					if c > 0 {
						add = true
						break
					}
				}
			}

			if add {
				total = append(total, counts)
			}
		}
	}

	return total, nil
}
