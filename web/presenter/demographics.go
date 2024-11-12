package presenter

import (
	"github.com/louisbranch/edulab"
	"golang.org/x/text/message"
)

type Demographic struct {
	edulab.Demographic
	printer *message.Printer
}

func (d Demographic) Prompt() string {
	if d.Demographic.Text != "" {
		return d.Demographic.Text
	}

	if d.Demographic.I18nKey != "" {
		return d.printer.Sprint(d.Demographic.I18nKey)
	}

	return d.printer.Sprint("Unknown Demographic Prompt")
}

func NewDemographic(d edulab.Demographic, printer *message.Printer) Demographic {
	return Demographic{
		Demographic: d,
		printer:     printer,
	}
}

func NewDemographics(ds []edulab.Demographic, printer *message.Printer) []Demographic {
	result := make([]Demographic, len(ds))
	for i, d := range ds {
		result[i] = NewDemographic(d, printer)
	}
	return result
}
