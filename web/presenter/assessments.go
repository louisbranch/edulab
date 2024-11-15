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
	return AssessmentType(a.printer, a.Assessment.Type)
}

func AssessmentType(printer *message.Printer, t edulab.AssessmentType) string {
	switch t {
	case edulab.AssessmentTypePre:
		return printer.Sprintf("Pre-Assessment")
	case edulab.AssessmentTypePost:
		return printer.Sprintf("Post-Assessment")
	default:
		return printer.Sprintf("Unknown Assessment Type")
	}
}

func AssessmentTypes(printer *message.Printer) [][]string {
	return [][]string{
		{string(edulab.AssessmentTypePre), AssessmentType(printer, edulab.AssessmentTypePre)},
		{string(edulab.AssessmentTypePost), AssessmentType(printer, edulab.AssessmentTypePost)},
	}
}
