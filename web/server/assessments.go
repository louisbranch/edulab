package server

import (
	"html/template"
	"net/http"

	"github.com/louisbranch/edulab"
	"github.com/louisbranch/edulab/web/presenter"
)

func (srv *Server) assessmentsHandler(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment, segments []string) {
	if len(segments) < 1 {
		srv.renderNotFound(w, r)
		return
	}

	pid := segments[0]
	assessment, err := srv.DB.FindAssessment(experiment.ID, pid)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	if len(segments) == 1 {
		srv.showAssessment(w, r, experiment, assessment)
		return
	}

	srv.renderNotFound(w, r)
}

func (srv *Server) showAssessment(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment, assessment edulab.Assessment) {

	printer, page := srv.i18n(w, r)

	page.Title = printer.Sprintf("Assessment: %s", assessment.Name)
	page.Partials = []string{"assessment"}
	page.Content = struct {
		Breadcrumbs   template.HTML
		Experiment    edulab.Experiment
		Assessment    edulab.Assessment
		Questions     []edulab.AssessmentQuestion
		QuestionTypes []presenter.QuestionType
		Texts         interface{}
	}{
		Breadcrumbs:   presenter.RenderBreadcrumbs(presenter.ExperimentsBreadcrumb(&experiment, printer)),
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
			Questions          string
			NewQuestion        string
			NoQuestions        string
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
			Questions:   printer.Sprintf("Questions"),
			NewQuestion: printer.Sprintf("New Question"),
			NoQuestions: printer.Sprintf("No questions yet"),
		},
	}

	srv.render(w, page)
}
