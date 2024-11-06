package presenter

import (
	"github.com/louisbranch/edulab"
	"golang.org/x/text/message"
)

type Experiment struct {
	edulab.Experiment
}

func ExperimentsList(sessions []edulab.Experiment, printer *message.Printer) []Experiment {
	list := []Experiment{}

	return list
}
