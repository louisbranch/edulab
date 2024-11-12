package presenter

import (
	"golang.org/x/text/message"

	"github.com/louisbranch/edulab"
)

type Question struct {
	edulab.Question
	Choices []edulab.QuestionChoice
}

type QuestionType struct {
	Value string
	Text  string
}

func QuestionTypes(printer *message.Printer) []QuestionType {
	return []QuestionType{
		{Value: string(edulab.SingleChoice), Text: printer.Sprintf("Single choice")},
		{Value: string(edulab.MultipleChoice), Text: printer.Sprintf("Multiple choice")},
		{Value: string(edulab.Text), Text: printer.Sprintf("Text")},
	}
}

func SortQuestions(questions []edulab.Question, qc map[string]Question) []Question {
	var sorted []Question
	for _, q := range questions {
		sorted = append(sorted, qc[q.ID])
	}
	return sorted
}
