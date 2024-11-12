package presenter

import (
	"github.com/louisbranch/edulab"
	"golang.org/x/text/message"
)

type Assessment struct {
	edulab.Assessment
	printer *message.Printer
}

func NewAssessment(assessment edulab.Assessment, printer *message.Printer) Assessment {
	return Assessment{
		Assessment: assessment,
		printer:    printer,
	}
}

func NewAssessments(assessments []edulab.Assessment, printer *message.Printer) []Assessment {
	result := make([]Assessment, len(assessments))
	for i, a := range assessments {
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
