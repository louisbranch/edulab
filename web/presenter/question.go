package presenter

import (
	"golang.org/x/text/message"

	"github.com/louisbranch/edulab"
)

type QuestionType struct {
	Value string
	Text  string
}

func QuestionTypes(printer *message.Printer) []QuestionType {
	return []QuestionType{
		{Value: string(edulab.SingleChoice), Text: printer.Sprintf("Single choice")},
		{Value: string(edulab.MultipleChoice), Text: printer.Sprintf("Multiple choice")},
		{Value: string(edulab.FreeForm), Text: printer.Sprintf("Free form")},
	}
}
