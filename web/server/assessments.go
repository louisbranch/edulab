package server

import (
	"html/template"
	"log"
	"net/http"

	"github.com/louisbranch/edulab"
	"github.com/louisbranch/edulab/web/presenter"
)

func (srv *Server) assessmentsHandler(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment, segments []string) {

	log.Print("[DEBUG] web/server/assessments.go: handling assessments")

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

func (srv *Server) showAssessment(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment, assessment edulab.Assessment) {

	questions, err := srv.DB.FindQuestions(assessment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	printer, page := srv.i18n(w, r)

	page.Title = printer.Sprint("Assessment")
	page.Partials = []string{"assessment"}
	page.Content = struct {
		Breadcrumbs template.HTML
		Experiment  edulab.Experiment
		Assessment  presenter.Assessment
		Questions   []edulab.Question
		Texts       interface{}
	}{
		Breadcrumbs: presenter.ExperimentBreadcrumb(experiment, printer),
		Experiment:  experiment,
		Assessment:  presenter.NewAssessment(assessment, printer),
		Questions:   questions,
		Texts: struct {
			Questions   string
			Prompt      string
			Actions     string
			Edit        string
			AddQuestion string
			NoQuestions string
			Preview     string
		}{
			Questions:   printer.Sprintf("Questions"),
			Prompt:      printer.Sprintf("Prompt"),
			Actions:     printer.Sprintf("Actions"),
			Edit:        printer.Sprintf("Edit"),
			AddQuestion: printer.Sprintf("Add Question"),
			NoQuestions: printer.Sprintf("No questions yet"),
			Preview:     printer.Sprintf("Preview"),
		},
	}

	srv.render(w, page)
}

func (srv *Server) previewAssessment(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment, assessment edulab.Assessment) {

	questions, err := srv.DB.FindQuestions(assessment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	qp := make(map[string]presenter.Question)
	for _, q := range questions {

		qp[q.ID] = presenter.Question{
			Question: q,
		}
	}

	choices, err := srv.DB.FindQuestionChoices(assessment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	for _, c := range choices {
		q, ok := qp[c.QuestionID]
		if !ok {
			log.Printf("[ERROR] web/server/assessments.go: question not found: %s", c.QuestionID)
			continue
		}
		q.Choices = append(q.Choices, c)
		qp[c.QuestionID] = q
	}

	printer, page := srv.i18n(w, r)

	page.Title = printer.Sprint("Preview Assessment")
	page.Partials = []string{"preview_assessment"}
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
		Questions:   presenter.SortQuestions(questions, qp),
		Texts: struct {
			Questions string
		}{
			Questions: printer.Sprintf("Questions"),
		},
	}

	srv.render(w, page)

}
