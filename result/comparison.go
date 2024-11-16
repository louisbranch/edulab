package result

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/louisbranch/edulab/stats"
)

type Comparison struct {
	headers []string
	data    map[string][]float64 // header -> scores
	rows    int
}

type AssessmentQuestions struct {
	AssessmentID string
	QuestionID   string
}

// NewComparison initializes a Comparison struct with scores across specified cohorts and assessments.
func NewComparison(r *Result, assessmentQuestions []AssessmentQuestions,
	cohorts []string) (*Comparison, error) {
	c := &Comparison{
		headers: []string{},
		data:    make(map[string][]float64),
	}

	cohortLabels := []string{"base", "intervention"}
	if len(cohorts) != 2 {
		cohortLabels = []string{}
		for _, cohortID := range cohorts {
			cohort := r.cohorts[cohortID]
			cohortLabels = append(cohortLabels, strings.ToLower(cohort.Name))
		}
	}

	// Populate score for each assignment question
	for _, val := range assessmentQuestions {

		assessement := r.assessments[val.AssessmentID]

		questionID := val.QuestionID
		question := r.questions[questionID]

		log.Printf("[DEBUG] Comparing question %s: %q\n", questionID, question.Text)

		scores, err := r.QuestionScore(questionID)
		if err != nil {
			return nil, err
		}

		for i, cohortID := range cohorts {
			label := cohortLabels[i]

			header := fmt.Sprintf("%s_%s", assessement.Type, label)
			c.headers = append(c.headers, header)

			score := scores[cohortID]
			if _, ok := c.data[header]; ok {
				return nil, fmt.Errorf("%s already exists", header)
			}

			c.data[header] = score

			n := len(score)
			if n > c.rows {
				c.rows = n
			}
		}
	}

	return c, nil
}

// ToCSV exports the comparison data to a CSV file.
func (c *Comparison) ToCSV(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers
	if err := writer.Write(c.headers); err != nil {
		return err
	}

	// Write rows
	for i := 0; i < c.rows; i++ {
		records := make([]string, len(c.headers))
		for j, header := range c.headers {
			scores := c.data[header]
			if len(scores) < i {
				records[j] = ""
			} else {
				records[j] = strconv.FormatFloat(scores[i], 'f', 2, 64)
			}
		}

		writer.Write(records)
	}

	return nil
}

func (c *Comparison) ToStatsData() []stats.Data {
	var data []stats.Data
	for i := 0; i < c.rows; i++ {
		preBase := c.data["pre_base"][i]
		postBase := c.data["post_base"][i]
		preIntervention := c.data["pre_intervention"][i]
		postIntervention := c.data["post_intervention"][i]

		data = append(data, stats.Data{
			PreBase:          preBase,
			PostBase:         postBase,
			PreIntervention:  preIntervention,
			PostIntervention: postIntervention,
		})
	}
	return data
}
