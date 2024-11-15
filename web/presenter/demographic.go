package presenter

import (
	"encoding/json"
	"fmt"

	"github.com/louisbranch/edulab"
)

type Demographics []Demographic

func (d Demographics) Labels() [][]string {
	var labels [][]string
	for _, dd := range d {
		var dl []string
		for _, o := range dd.Options {
			dl = append(dl, o.Text)
		}
		labels = append(labels, dl)
	}
	return labels
}

func (d Demographics) Values(participations []edulab.Participation) ([][]int, error) {

	// Filter and cast to map[string][]string
	result := make(map[string][]string)
	for _, p := range participations {
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

	for _, dd := range d {
		dv := make([]int, len(dd.Options))
		for i, o := range dd.Options {
			for _, v := range result[dd.ID] {
				if v == o.ID {
					dv[i]++
				}
			}
		}

		values = append(values, dv)
	}

	return values, nil
}

type Demographic struct {
	edulab.Demographic
	Options []edulab.DemographicOption
}

func NewDemographics(ds []edulab.Demographic, dos []edulab.DemographicOption) Demographics {
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
