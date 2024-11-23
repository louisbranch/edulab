package result

import (
	"reflect"
	"testing"

	"github.com/louisbranch/edulab"
	"github.com/louisbranch/edulab/mock"
)

func TestCountChoicesByCohorts(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := mock.NewDB()

		err := db.CreateParticipation(&edulab.Participation{
			ExperimentID:  "1",
			AssessmentID:  "1",
			CohortID:      "1",
			ParticipantID: "1",
			Answers:       []byte(`{"1":["1"],"2":["4","5"],"3":["Words..."]}`),
		})
		if err != nil {
			t.Errorf("UpdateParticipation() error = %v, want nil", err)
			return
		}

		err = db.CreateParticipation(&edulab.Participation{
			ExperimentID:  "1",
			AssessmentID:  "2",
			CohortID:      "1",
			ParticipantID: "1",
			Answers:       []byte(`{"4":["9"],"5":["12","13"],"6":["Words..."]}`),
		})
		if err != nil {
			t.Errorf("UpdateParticipation() error = %v, want nil", err)
			return
		}

		experiment := edulab.Experiment{
			ID: "1",
		}

		actual, err := CountChoicesByCohorts(db, experiment)
		if err != nil {
			t.Errorf("CountChoicesByCohorts() error = %v, want nil", err)
			return
		}

		expected := [][][]int{
			{{1, 0, 0}, {0, 0, 0}},       // Question 1
			{{1, 1, 0, 0}, {0, 0, 0, 0}}, // Question 2
			{{0, 1, 0}, {0, 0, 0}},       // Question 4
			{{0, 1, 1, 0}, {0, 0, 0, 0}}, // Question 5
		}

		// Compare using reflect.DeepEqual
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("CountChoicesByCohorts() = %v, want %v", actual, expected)
		}

	})
}
