package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/louisbranch/edulab"
	"github.com/louisbranch/edulab/web/presenter"
)

func (srv *Server) questionsHandler(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment, assessment edulab.Assessment, segments []string) {

	log.Print("[DEBUG] web/server/questions.go: handling questions")

	if len(segments) < 1 && r.Method == http.MethodPost {
		srv.createQuestion(w, r, experiment, assessment)
		return
	}

	pid := segments[0]

	switch pid {
	case "new":
		srv.newQuestionForm(w, r, experiment, assessment)
		return
	default:
		srv.showQuestion(w, r, experiment, assessment, pid)
		return
	}
}

func (srv *Server) newQuestionForm(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment, assessment edulab.Assessment) {

	printer, page := srv.i18n(w, r)

	page.Title = printer.Sprintf("New Question")
	page.Partials = []string{"new_question"}
	page.Content = struct {
		Breadcrumbs   template.HTML
		Experiment    edulab.Experiment
		Assessment    edulab.Assessment
		QuestionTypes []presenter.QuestionType
		Texts         interface{}
	}{
		Breadcrumbs:   presenter.AssessmentBreadcrumb(experiment, assessment, printer),
		Experiment:    experiment,
		Assessment:    assessment,
		QuestionTypes: presenter.QuestionTypes(printer),
		Texts: struct {
			Prompt             string
			PromptHelp         string
			PromptPlaceholder  string
			Type               string
			Choices            string
			ChoicesHelp        string
			ChoicePlaceholders []string
			Correct            string
			Create             string
			NewQuestion        string
		}{
			Prompt:            printer.Sprintf("Prompt"),
			PromptHelp:        printer.Sprintf("Markdown supported"),
			PromptPlaceholder: printer.Sprintf("e.g. What is the best explanation for the cause of Earth's seasons?"),
			Type:              printer.Sprintf("Type"),
			Choices:           printer.Sprintf("Choices"),
			ChoicesHelp:       printer.Sprintf("Markdown supported. Empty choices will be ignored."),
			ChoicePlaceholders: []string{
				printer.Sprintf("e.g. The tilt of Earth's axis"),
				printer.Sprintf("e.g. The distance from the Sun"),
				printer.Sprintf("e.g. The Earth's elliptical orbit"),
				printer.Sprintf("e.g. The Earth's rotation"),
				printer.Sprintf("e.g. The Earth's revolution"),
			},
			Correct:     printer.Sprintf("Correct"),
			Create:      printer.Sprintf("Create"),
			NewQuestion: printer.Sprintf("New Question"),
		},
	}

	srv.render(w, page)
}

func (srv *Server) showQuestion(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment, assessment edulab.Assessment, pid string) {

	printer, page := srv.i18n(w, r)

	question, err := srv.DB.FindQuestion(assessment.ID, pid)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	page.Title = printer.Sprintf("Question: %s", question.Prompt)
	page.Partials = []string{"question"}
	page.Content = struct {
		Breadcrumbs template.HTML
		Experiment  edulab.Experiment
		Assessment  edulab.Assessment
		Question    edulab.Question
	}{
		Experiment: experiment,
		Assessment: assessment,
		Question:   question,
	}

	srv.render(w, page)
}

func (srv *Server) createQuestion(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment, assessment edulab.Assessment) {

	err := r.ParseForm()
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	prompt := r.FormValue("prompt")
	qtype := r.FormValue("type")
	choices := r.Form["choices[]"]

	correct := make([]bool, len(choices))

	// Parse `correct[]` values as indices
	checkedIndices := r.Form["correct[]"]
	for _, indexStr := range checkedIndices {
		index, err := strconv.Atoi(indexStr)
		if err == nil && index < len(choices) {
			correct[index] = true
		}
	}

	question := edulab.Question{
		AssessmentID: assessment.ID,
		Prompt:       prompt,
		Type:         edulab.InputType(qtype),
	}

	err = srv.DB.CreateQuestion(&question)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	for i, choice := range choices {
		if strings.Trim(choice, " ") == "" {
			continue
		}

		qc := edulab.QuestionChoice{
			QuestionID: question.ID,
			Text:       choice,
			IsCorrect:  correct[i],
		}

		err = srv.DB.CreateQuestionChoice(&qc)
		if err != nil {
			srv.renderError(w, r, err)
			return
		}
	}

	uri := fmt.Sprintf("/edulab/experiments/%s/assessments/%s/", experiment.PublicID, assessment.PublicID)
	http.Redirect(w, r, uri, http.StatusFound)
}
