package presenter

import (
	"log"

	"golang.org/x/text/message"

	"github.com/louisbranch/edulab"
)

type Question struct {
	edulab.Question
	Choices []edulab.QuestionChoice `json:"choices"`
}

type QuestionType struct {
	Value string
	Text  string
}

func GroupQuestions(questions []edulab.Question, choices []edulab.QuestionChoice) []Question {

	dict := make(map[string]Question)
	for _, q := range questions {

		dict[q.ID] = Question{
			Question: q,
		}
	}

	for _, c := range choices {
		q, ok := dict[c.QuestionID]
		if !ok {
			log.Printf("[WARN] Choice without question: %s -> %s", c.QuestionID, c.ID)
			continue
		}
		q.Choices = append(q.Choices, c)
		dict[c.QuestionID] = q
	}

	var sorted []Question
	for _, q := range questions {
		sorted = append(sorted, dict[q.ID])
	}
	return sorted
}

func QuestionTypes(printer *message.Printer) []QuestionType {
	return []QuestionType{
		{Value: string(edulab.InputSingle), Text: printer.Sprintf("Single Choice")},
		{Value: string(edulab.InputMultiple), Text: printer.Sprintf("Multiple Choice")},
		{Value: string(edulab.InputText), Text: printer.Sprintf("Text")},
	}
}
