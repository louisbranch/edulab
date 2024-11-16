package result

import (
	"sort"
	"strconv"

	"github.com/louisbranch/edulab"
)

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

func (r *Result) Valid() bool {
	return len(r.cohorts) > 1 && len(r.assessments) > 1
}

func (r *Result) Participation() bool {
	return len(r.participation) > 0
}

func (r *Result) ComparisonPairs() ([]string, [][]AssessmentQuestions) {

	cohortIDs := make([]string, 0, len(r.cohorts))
	for id := range r.cohorts {
		cohortIDs = append(cohortIDs, id)
	}

	sort.Slice(cohortIDs, func(i, j int) bool {
		return cohortIDs[i] < cohortIDs[j]
	})

	questions := make(map[string][]AssessmentQuestions)

	for _, q := range r.questions {
		mq := questions[q.Text]
		a := AssessmentQuestions{
			AssessmentID: q.AssessmentID,
			QuestionID:   q.ID,
		}

		mq = append(mq, a)
		questions[q.Text] = mq
	}

	var items [][]AssessmentQuestions
	for _, mq := range questions {
		if len(mq) < 2 {
			continue
		}

		sort.Slice(mq, func(i, j int) bool {
			a1, _ := strconv.Atoi(mq[i].AssessmentID)
			a2, _ := strconv.Atoi(mq[i].AssessmentID)

			if a1 != a2 {
				return a1 < a2
			}

			q1, _ := strconv.Atoi(mq[i].QuestionID)
			q2, _ := strconv.Atoi(mq[j].QuestionID)

			return q1 < q2
		})

		items = append(items, mq)
	}

	sort.Slice(items, func(i, j int) bool {
		a1, a2 := items[i][0].AssessmentID, items[j][0].AssessmentID
		q1, q2 := items[i][0].QuestionID, items[j][0].QuestionID
		return a1 < a2 || (a1 == a2 && q1 < q2)
	})

	return cohortIDs, items
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
