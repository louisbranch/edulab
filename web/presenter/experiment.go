package presenter

import (
	"math"
	"time"

	"golang.org/x/text/message"

	"github.com/louisbranch/edulab"
)

type Experiment struct {
	edulab.Experiment
	ElapsedTime  string
	Participants int
}

func ExperimentsList(experiments []edulab.Experiment, printer *message.Printer) []Experiment {
	list := []Experiment{}

	for _, experiment := range experiments {
		var elapsed string

		mins := time.Since(experiment.CreatedAt).Minutes()
		hours := time.Since(experiment.CreatedAt).Hours()

		switch {
		case mins < 1:
			elapsed = printer.Sprintf("Less than one min ago")
		case mins < 90:
			min := int64(math.Ceil(mins))
			elapsed = printer.Sprintf("%d mins ago", min)
		case hours < 24:
			hr := int64(math.Ceil(hours))
			elapsed = printer.Sprintf("%d hours ago", hr)
		default:
			days := int64(math.Ceil(hours / 24))
			elapsed = printer.Sprintf("%d days ago", days)
		}

		list = append(list, Experiment{experiment, elapsed, 0})
	}

	return list
}
