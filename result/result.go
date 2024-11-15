package result

import "github.com/louisbranch/edulab"

type Result struct {
	db            edulab.Database
	experimentID  string
	participants  map[string]edulab.Participant
	assessments   map[string]edulab.Assessment
	cohorts       map[string]edulab.Cohort
	questions     map[string]edulab.Question
	choices       map[string][]edulab.QuestionChoice
	participation map[string][]edulab.Participation // Map from participantID to their Participation records
}

// New initializes a new Result instance, loading data into memory.
func New(db edulab.Database, experimentID string) (*Result, error) {
	r := &Result{
		db:            db,
		experimentID:  experimentID,
		participants:  make(map[string]edulab.Participant),
		assessments:   make(map[string]edulab.Assessment),
		cohorts:       make(map[string]edulab.Cohort),
		questions:     make(map[string]edulab.Question),
		choices:       make(map[string][]edulab.QuestionChoice),
		participation: make(map[string][]edulab.Participation),
	}

	// Load participants
	participants, err := db.FindParticipants(experimentID)
	if err != nil {
		return nil, err
	}
	for _, p := range participants {
		r.participants[p.ID] = p
	}

	// Load cohorts
	cohorts, err := db.FindCohorts(experimentID)
	if err != nil {
		return nil, err
	}
	for _, c := range cohorts {
		r.cohorts[c.ID] = c
	}

	// Load questions
	assessments, err := db.FindAssessments(experimentID)
	if err != nil {
		return nil, err
	}
	for _, a := range assessments {
		r.assessments[a.ID] = a

		questions, err := db.FindQuestions(a.ID)
		if err != nil {
			return nil, err
		}
		allChoices, err := db.FindQuestionChoices(a.ID)
		if err != nil {
			return nil, err
		}
		for _, q := range questions {
			var choices []edulab.QuestionChoice
			for _, choice := range allChoices {
				if choice.QuestionID == q.ID {
					choices = append(choices, choice)
				}
			}
			r.questions[q.ID] = q
			r.choices[q.ID] = choices
		}
	}

	// Load participation data
	for _, p := range participants {
		participations, err := db.FindParticipationsByParticipant(experimentID, p.ID)
		if err != nil {
			return nil, err
		}
		r.participation[p.ID] = participations
	}

	return r, nil
}
