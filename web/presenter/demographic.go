package presenter

import (
	"github.com/louisbranch/edulab"
)

type Demographic struct {
	edulab.Demographic
	Options []edulab.DemographicOption
}

func SortDemographics(ds []edulab.Demographic, do map[string]Demographic) []Demographic {
	var sorted []Demographic
	for _, d := range ds {
		sorted = append(sorted, do[d.ID])
	}
	return sorted
}
