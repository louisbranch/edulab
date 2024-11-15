package presenter

import (
	"github.com/louisbranch/edulab"
)

type Demographic struct {
	edulab.Demographic
	Options []edulab.DemographicOption
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
