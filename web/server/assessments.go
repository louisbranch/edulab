package server

import (
	"html/template"
	"log"
	"net/http"

	"github.com/louisbranch/edulab"
	"github.com/louisbranch/edulab/web/presenter"
)

// assessmentsHandler handles the assessments subroutes in an experiment for the instructor.
func (srv *Server) assessmentsHandler(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment, segments []string) {

	log.Print("[DEBUG] web/server/assessments.go: handling assessments")

	if len(segments) < 1 {
		srv.listAssessments(w, r, experiment)
		return
	}

	pid := segments[0]
	assessment, err := srv.DB.FindAssessment(experiment.ID, pid)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	if len(segments) == 1 {
		srv.editAssessment(w, r, experiment, assessment)
		return
	}

	switch segments[1] {
	case "preview":
		srv.previewAssessment(w, r, experiment, assessment)
		return
	case "questions":
		srv.questionsHandler(w, r, experiment, assessment, segments[2:])
		return
	default:
		srv.renderNotFound(w, r)
		return
	}
}

// listAssessments lists the assessments for an experiment to the instructor.
func (srv *Server) listAssessments(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment) {

	assessments, err := srv.DB.FindAssessments(experiment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	printer, page := srv.i18n(w, r)

	page.Title = printer.Sprintf("Assessments")
	page.Partials = []string{"assessments"}
	page.Content = struct {
		Breadcrumbs template.HTML
		Experiment  edulab.Experiment
		Assessments []presenter.Assessment
		Texts       interface{}
	}{
		Breadcrumbs: presenter.ExperimentBreadcrumb(experiment, printer),
		Experiment:  experiment,
		Assessments: presenter.NewAssessments(assessments, printer),
		Texts: struct {
			Title      string
			Type       string
			Questions  string
			Actions    string
			Edit       string
			Add        string
			Empty      string
			Preview    string
			ComingSoon string
		}{
			Title:      printer.Sprintf("Assessments"),
			Type:       printer.Sprintf("Type"),
			Questions:  printer.Sprintf("Questions"),
			Actions:    printer.Sprintf("Actions"),
			Add:        printer.Sprintf("Add Assessment"),
			Empty:      printer.Sprintf("No assessments yet"),
			Edit:       printer.Sprintf("Edit"),
			Preview:    printer.Sprintf("Preview"),
			ComingSoon: printer.Sprintf("Coming Soon"),
		},
	}

	srv.render(w, page)
}

// editAssessment displays the assessment editor to the instructor.
func (srv *Server) editAssessment(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment, assessment edulab.Assessment) {

	questions, err := srv.DB.FindQuestions(assessment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	printer, page := srv.i18n(w, r)

	page.Title = printer.Sprintf("Assessment")
	page.Partials = []string{"assessment_edit"}
	page.Content = struct {
		Breadcrumbs template.HTML
		Experiment  edulab.Experiment
		Assessment  presenter.Assessment
		Questions   []edulab.Question
		Texts       interface{}
	}{
		Breadcrumbs: presenter.AssessmentsBreadcrumb(experiment, printer),
		Experiment:  experiment,
		Assessment:  presenter.NewAssessment(assessment, printer),
		Questions:   questions,
		Texts: struct {
			Description            string
			DescriptionHelp        string
			DescriptionPlaceholder string
			Questions              string
			Text                   string
			Actions                string
			Edit                   string
			Save                   string
			Add                    string
			Empty                  string
			Preview                string
			ComingSoon             string
		}{
			Description:            printer.Sprintf("Description"),
			DescriptionHelp:        printer.Sprintf("Optional. Markdown supported."),
			DescriptionPlaceholder: printer.Sprintf("e.g. Gauge your current knowledge about the causes of Earth's..."),
			Questions:              printer.Sprintf("Questions"),
			Text:                   printer.Sprintf("Text"),
			Actions:                printer.Sprintf("Actions"),
			Edit:                   printer.Sprintf("Edit"),
			Save:                   printer.Sprintf("Save"),
			Add:                    printer.Sprintf("Add Question"),
			Empty:                  printer.Sprintf("No questions yet"),
			Preview:                printer.Sprintf("Preview Assessment"),
			ComingSoon:             printer.Sprintf("Coming Soon"),
		},
	}

	srv.render(w, page)
}

// previewAssessment displays the assessment preview to the instructor.
func (srv *Server) previewAssessment(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment, assessment edulab.Assessment) {

	questions, err := srv.DB.FindQuestions(assessment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	choices, err := srv.DB.FindQuestionChoices(assessment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	qp := presenter.GroupQuestions(questions, choices)

	printer, page := srv.i18n(w, r)

	page.Title = printer.Sprintf("Preview Assessment")
	page.Partials = []string{"assessment_preview"}
	page.Content = struct {
		Breadcrumbs template.HTML
		Experiment  edulab.Experiment
		Assessment  presenter.Assessment
		Questions   []presenter.Question
		Texts       interface{}
	}{
		Breadcrumbs: presenter.AssessmentBreadcrumb(experiment, assessment, printer),
		Experiment:  experiment,
		Assessment:  presenter.NewAssessment(assessment, printer),
		Questions:   qp,
		Texts: struct {
			Questions string
			Submit    string
			Back      string
		}{
			Questions: printer.Sprintf("Questions"),
			Submit:    printer.Sprintf("Submit"),
			Back:      printer.Sprintf("Back"),
		},
	}

	srv.render(w, page)

}

// showAssessment displays the assessment to the participant.
func (srv *Server) showAssessment(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment, cohort edulab.Cohort,
	participant edulab.Participant, assessment edulab.Assessment) {

	printer, page := srv.i18n(w, r)

	questions, err := srv.DB.FindQuestions(assessment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	choices, err := srv.DB.FindQuestionChoices(assessment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	qp := presenter.GroupQuestions(questions, choices)

	page.Title = printer.Sprintf("%s - %s", experiment.Name, assessment.Type)
	page.Partials = []string{"assessment_participate"}
	page.Content = struct {
		edulab.Experiment
		edulab.Cohort
		edulab.Participant
		presenter.Assessment
		Questions []presenter.Question
		Texts     interface{}
	}{
		Experiment:  experiment,
		Cohort:      cohort,
		Participant: participant,
		Assessment:  presenter.NewAssessment(assessment, printer),
		Questions:   qp,
		Texts: struct {
			Submit string
		}{
			Submit: printer.Sprintf("Submit"),
		},
	}

	srv.render(w, page)
}
