package presenter

import (
	"encoding/json"
	"fmt"

	"github.com/louisbranch/edulab"
)

type DemographicsResult struct {
	Demographics   []Demographic
	cohorts        []edulab.Cohort
	participants   []edulab.Participant
	participations []edulab.Participation
	Categories     [][]string
	Cohorts        []string
	Data           [][][]int // [category][cohort][count]
}

func (dr DemographicsResult) categories() [][]string {
	var labels [][]string
	for _, d := range dr.Demographics {
		var dl []string
		for _, o := range d.Options {
			dl = append(dl, o.Text)
		}
		labels = append(labels, dl)
	}
	return labels
}

func (dr DemographicsResult) cohortNames() []string {
	var names []string
	for _, c := range dr.cohorts {
		names = append(names, c.Name)
	}
	return names
}

func NewDemographicsResult(ds []edulab.Demographic, dos []edulab.DemographicOption,
	cohorts []edulab.Cohort, participants []edulab.Participant,
	participations []edulab.Participation) (DemographicsResult, error) {

	dr := DemographicsResult{
		Demographics:   NewDemographics(ds, dos),
		cohorts:        cohorts,
		participants:   participants,
		participations: participations,
	}
	data, err := dr.data()
	if err != nil {
		return dr, err
	}

	dr.Categories = dr.categories()
	dr.Data = data
	dr.Cohorts = dr.cohortNames()

	return dr, nil
}

func (dr DemographicsResult) data() ([][][]int, error) {

	// Map of participant ID to cohort ID
	participants := make(map[string]string)
	for _, p := range dr.participants {
		participants[p.ID] = p.CohortID
	}

	// cohort ID -> demographic ID -> option ID -> count
	data := make(map[string]map[string]map[string]int)

	for _, p := range dr.participations {
		if p.Demographics == nil {
			continue
		}

		// Unmarshal into a map of interface{}
		var tempMap map[string]interface{}
		if err := json.Unmarshal(p.Demographics, &tempMap); err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return nil, err
		}

		for demographicID, values := range tempMap {
			// Check if the value is a slice of strings
			var stringArray []string
			if str, ok := values.(string); ok {
				stringArray = []string{str}
			} else if array, ok := values.([]interface{}); ok {
				for _, item := range array {
					if str, isString := item.(string); isString {
						stringArray = append(stringArray, str)
					} else {
						stringArray = nil // Skip if any item isn't a string
						break
					}
				}
			}
			// Add to result only if all items were strings
			if stringArray != nil {
				cohortID := participants[p.ParticipantID]
				if _, ok := data[cohortID]; !ok {
					data[cohortID] = make(map[string]map[string]int)
				}
				if _, ok := data[cohortID][demographicID]; !ok {
					data[cohortID][demographicID] = make(map[string]int)
				}

				for _, optionID := range stringArray {
					data[cohortID][demographicID][optionID]++
				}
			}
		}
	}

	values := make([][][]int, len(dr.Demographics))

	for i, d := range dr.Demographics {
		dv := make([][]int, len(dr.cohorts))

		for j, c := range dr.cohorts {
			do := make([]int, len(d.Options))
			for k, o := range d.Options {
				count := data[c.ID][d.ID][o.ID]
				do[k] = count
			}

			dv[j] = do
		}

		values[i] = dv
	}

	return values, nil
}
