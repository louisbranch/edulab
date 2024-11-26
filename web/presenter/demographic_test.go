package presenter

import (
	"reflect"
	"testing"

	"github.com/louisbranch/edulab"
)

func TestNewDemographics(t *testing.T) {
	ds := []edulab.Demographic{
		{ID: "1", Text: "Age"},
		{ID: "2", Text: "Gender"},
	}
	dos := []edulab.DemographicOption{
		{ID: "1", DemographicID: "1", Text: "18-24"},
		{ID: "2", DemographicID: "1", Text: "25-34"},
		{ID: "3", DemographicID: "2", Text: "Male"},
		{ID: "4", DemographicID: "2", Text: "Female"},
	}

	expected := []Demographic{
		{
			Demographic: ds[0],
			Options: []edulab.DemographicOption{
				dos[0], dos[1],
			},
		},
		{
			Demographic: ds[1],
			Options: []edulab.DemographicOption{
				dos[2], dos[3],
			},
		},
	}

	result := NewDemographics(ds, dos)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("NewDemographics() = %v, want %v", result, expected)
	}
}

func TestSortDemographics(t *testing.T) {
	ds := []edulab.Demographic{
		{ID: "1", Text: "Age"},
		{ID: "2", Text: "Gender"},
	}
	do := map[string]Demographic{
		"1": {Demographic: edulab.Demographic{ID: "1", Text: "Age"}},
		"2": {Demographic: edulab.Demographic{ID: "2", Text: "Gender"}},
	}

	expected := []Demographic{
		{Demographic: edulab.Demographic{ID: "1", Text: "Age"}},
		{Demographic: edulab.Demographic{ID: "2", Text: "Gender"}},
	}

	result := SortDemographics(ds, do)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("SortDemographics() = %v, want %v", result, expected)
	}
}
