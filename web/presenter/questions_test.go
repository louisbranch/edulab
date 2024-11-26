package presenter

import (
	"testing"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/louisbranch/edulab"
)

func TestGroupQuestions(t *testing.T) {
	questions := []edulab.Question{
		{ID: "q1", Text: "Question 1"},
		{ID: "q2", Text: "Question 2"},
	}

	choices := []edulab.QuestionChoice{
		{ID: "c1", QuestionID: "q1", Text: "Choice 1"},
		{ID: "c2", QuestionID: "q1", Text: "Choice 2"},
		{ID: "c3", QuestionID: "q2", Text: "Choice 3"},
	}

	grouped := GroupQuestions(questions, choices)

	if len(grouped) != 2 {
		t.Fatalf("expected 2 grouped questions, got %d", len(grouped))
	}

	if len(grouped[0].Choices) != 2 {
		t.Errorf("expected 2 choices for question 1, got %d", len(grouped[0].Choices))
	}

	if len(grouped[1].Choices) != 1 {
		t.Errorf("expected 1 choice for question 2, got %d", len(grouped[1].Choices))
	}
}

func TestQuestionTypes(t *testing.T) {
	printer := message.NewPrinter(language.English)
	types := QuestionTypes(printer)

	expected := []QuestionType{
		{Value: string(edulab.InputSingle), Text: "Single Choice"},
		{Value: string(edulab.InputMultiple), Text: "Multiple Choice"},
		{Value: string(edulab.InputText), Text: "Text"},
	}

	if len(types) != len(expected) {
		t.Fatalf("expected %d question types, got %d", len(expected), len(types))
	}

	for i, typ := range types {
		if typ != expected[i] {
			t.Errorf("expected type %v, got %v", expected[i], typ)
		}
	}
}
