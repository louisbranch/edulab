package presenter

import (
	"golang.org/x/text/message"

	"github.com/louisbranch/edulab"
)

type Assessment struct {
	edulab.Assessment
	printer *message.Printer
}

func NewAssessment(a edulab.Assessment, printer *message.Printer) Assessment {
	return Assessment{
		Assessment: a,
		printer:    printer,
	}
}

func NewAssessments(as []edulab.Assessment, printer *message.Printer) []Assessment {
	result := make([]Assessment, len(as))
	for i, a := range as {
		result[i] = NewAssessment(a, printer)
	}
	return result
}

func (a Assessment) Type() string {
	switch a.Assessment.Type {
	case edulab.AssessmentTypePre:
		return a.printer.Sprintf("Pre-Assessment")
	case edulab.AssessmentTypePos:
		return a.printer.Sprintf("Post-Assessment")
	default:
		return a.printer.Sprintf("Unknown Assessment Type")
	}
}
