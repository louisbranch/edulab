package presenter

import (
	"encoding/json"
	"fmt"

	"github.com/louisbranch/edulab"
)

type Demographic struct {
	edulab.Demographic
	Options []edulab.DemographicOption
}

type DemographicsResult struct {
	Demographics   []Demographic
	participants   []edulab.Participant
	participations []edulab.Participation
	Categories     [][]string
	Cohorts        []string
	//Data           [][][]int // [category][cohort][count]
	Data [][]int // [category][count]
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

func NewDemographicsResult(ds []edulab.Demographic, dos []edulab.DemographicOption,
	participants []edulab.Participant, participations []edulab.Participation) (DemographicsResult, error) {

	dr := DemographicsResult{
		Demographics:   NewDemographics(ds, dos),
		participants:   participants,
		participations: participations,
	}
	data, err := dr.data()
	if err != nil {
		return dr, err
	}

	dr.Categories = dr.categories()
	dr.Data = data

	return dr, nil
}

func (dr DemographicsResult) data() ([][]int, error) {

	// Filter and cast to map[string][]string
	result := make(map[string][]string)
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

		for k, v := range tempMap {
			// Check if the value is a slice of strings
			var stringArray []string
			if str, ok := v.(string); ok {
				stringArray = []string{str}
			} else if array, ok := v.([]interface{}); ok {
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
				if _, ok := result[k]; !ok {
					result[k] = []string{}
				}
				result[k] = append(result[k], stringArray...)
			}
		}
	}

	var values [][]int

	for _, d := range dr.Demographics {
		dv := make([]int, len(d.Options))
		for i, o := range d.Options {
			for _, v := range result[d.ID] {
				if v == o.ID {
					dv[i]++
				}
			}
		}

		values = append(values, dv)
	}

	return values, nil
}

func NewDemographics(ds []edulab.Demographic, dos []edulab.DemographicOption) []Demographic {
	var do []Demographic
	for _, d := range ds {
		nd := Demographic{
			Demographic: d,
		}
		for _, o := range dos {
			if o.DemographicID == d.ID {
				nd.Options = append(nd.Options, o)
			}
		}
		do = append(do, nd)
	}
	return do
}

func SortDemographics(ds []edulab.Demographic, do map[string]Demographic) []Demographic {
	var sorted []Demographic
	for _, d := range ds {
		sorted = append(sorted, do[d.ID])
	}
	return sorted
}
