package result

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Comparison struct {
	headers []string
	data    map[string]map[string]float64 // Participant ID -> [header -> score]
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
		data:    map[string]map[string]float64{},
	}

	// Create headers and populate scores for each participant
	for _, val := range assessmentQuestions {
		assessmentID := val.AssessmentID
		questionID := val.QuestionID

		assessment := r.assessments[assessmentID]
		question := r.questions[questionID]

		log.Printf("[DEBUG] Comparing question %s: %q\n", questionID, question.Text)

		for _, cohortID := range cohorts {
			cohort := r.cohorts[cohortID]

			header := fmt.Sprintf("%s_%s", assessment.Type, strings.ToLower(cohort.Name))
			c.headers = append(c.headers, header)

			// Get scores for the question
			scores, err := r.QuestionScore(questionID)
			if err != nil {
				return nil, err
			}

			// Add scores to the data map for each participant
			for participantID, score := range scores {
				if _, exists := c.data[participantID]; !exists {
					c.data[participantID] = map[string]float64{}
				}
				c.data[participantID][header] = score
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
	for _, row := range c.data {
		record := make([]string, len(c.headers))
		hasScores := false
		for i, header := range c.headers {
			if score, exists := row[header]; exists {
				record[i] = strconv.FormatFloat(score, 'f', 2, 64)
				hasScores = true
			} else {
				record[i] = "" // Leave blank if no score
			}
		}
		if hasScores {
			writer.Write(record)
		}
	}
	return nil
}

// ToJSON exports the comparison data to a JSON file.
func (c *Comparison) ToJSON(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// JSON data structure: map of participant scores by headers
	jsonData := map[string]map[string]float64{}
	for participantID, scores := range c.data {
		jsonData[participantID] = scores
	}

	encoder := json.NewEncoder(file)
	return encoder.Encode(jsonData)
}

// GetData returns the comparison data for direct use in Go.
func (c *Comparison) GetData() map[string]map[string]float64 {
	return c.data
}
