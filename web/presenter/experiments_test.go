package presenter

import (
	"testing"
	"time"

	"github.com/louisbranch/edulab"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func TestExperimentsList(t *testing.T) {
	printer := message.NewPrinter(language.English)

	tests := []struct {
		name        string
		experiments []edulab.Experiment
		expected    []Experiment
	}{
		{
			name: "less than one minute ago",
			experiments: []edulab.Experiment{
				{CreatedAt: time.Now().Add(-30 * time.Second)},
			},
			expected: []Experiment{
				{ElapsedTime: "Less than one min ago"},
			},
		},
		{
			name: "minutes ago",
			experiments: []edulab.Experiment{
				{CreatedAt: time.Now().Add(-9 * time.Minute)},
			},
			expected: []Experiment{
				{ElapsedTime: "10 mins ago"},
			},
		},
		{
			name: "hours ago",
			experiments: []edulab.Experiment{
				{CreatedAt: time.Now().Add(-110 * time.Minute)},
			},
			expected: []Experiment{
				{ElapsedTime: "2 hours ago"},
			},
		},
		{
			name: "days ago",
			experiments: []edulab.Experiment{
				{CreatedAt: time.Now().Add(-47 * time.Hour)},
			},
			expected: []Experiment{
				{ElapsedTime: "2 days ago"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExperimentsList(tt.experiments, printer)
			for i, experiment := range got {
				if experiment.ElapsedTime != tt.expected[i].ElapsedTime {
					t.Errorf("expected %s, got %s", tt.expected[i].ElapsedTime, experiment.ElapsedTime)
				}
			}
		})
	}
}
