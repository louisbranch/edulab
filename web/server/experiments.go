package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/louisbranch/edulab"
	"github.com/louisbranch/edulab/web/presenter"
	"github.com/louisbranch/edulab/wizard"
)

func (srv *Server) experimentsHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("[DEBUG] Routing experiments")

	path := r.URL.Path
	path = strings.TrimPrefix(path, "/experiments/")

	segments := strings.Split(path, "/")

	pid := segments[0]

	if pid == "" {
		switch r.Method {
		case "GET":
			srv.listExperiments(w, r)
		case "POST":
			srv.createExperiment(w, r)
		default:
			http.NotFound(w, r)
		}
		return
	}

	if pid == "new" {
		srv.newExperimentForm(w, r)
		return
	}

	experiment, err := srv.DB.FindExperiment(pid)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	if len(segments) == 1 && r.Method == "POST" {
		srv.updateExperiment(w, r, experiment)
		return
	}

	if len(segments) > 1 {
		switch segments[1] {
		case "edit":
			srv.editExperiment(w, r, experiment)
			return
		case "assessments":
			srv.assessmentsHandler(w, r, experiment, segments[2:])
			return
		case "demographics":
			srv.demographicsHandler(w, r, experiment, segments[2:])
			return
		case "cohorts":
			srv.cohortsHandler(w, r, experiment, segments[2:])
			return
		case "participate":
			srv.participateHandler(w, r, experiment)
			return
		case "results":
			srv.resultsHandler(w, r, experiment, segments[2:])
			return
		default:
			srv.renderNotFound(w, r)
			return
		}
	}

	srv.showExperiment(w, r, experiment)
}

func (srv *Server) newExperimentForm(w http.ResponseWriter, r *http.Request) {
	printer, page := srv.i18n(w, r)

	page.Title = printer.Sprintf("New Experiment")
	page.Partials = []string{"experiment_new"}
	page.Content = struct {
		Breadcrumbs template.HTML
		Texts       interface{}
	}{
		Breadcrumbs: presenter.HomeBreadcrumbs(printer),
		Texts: struct {
			Title                  string
			Name                   string
			NamePlaceholder        string
			Description            string
			DescriptionHelp        string
			DescriptionPlaceholder string
			Create                 string
		}{
			Title:                  printer.Sprintf("New Experiment"),
			Name:                   printer.Sprintf("Name"),
			NamePlaceholder:        printer.Sprintf("e.g. Earth's Seasons"),
			Description:            printer.Sprintf("Description"),
			DescriptionHelp:        printer.Sprintf("Optional. Not visible to participants."),
			DescriptionPlaceholder: printer.Sprintf("e.g. This experiment will compare 2 cohorts of students. One attending a traditional lecture and the other a workshop..."),
			Create:                 printer.Sprintf("Create"),
		},
	}

	srv.render(w, page)
}

func (srv *Server) createExperiment(w http.ResponseWriter, r *http.Request) {

	printer, _ := srv.i18n(w, r)

	err := r.ParseForm()
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	form := r.PostForm

	experiment := &edulab.Experiment{
		PublicID:    srv.newPublicID(2),
		Name:        form.Get("name"),
		Description: form.Get("description"),
	}

	err = srv.DB.CreateExperiment(experiment)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	err = wizard.Demographics(srv.DB, *experiment)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	types := []edulab.AssessmentType{
		edulab.AssessmentTypePre,
		edulab.AssessmentTypePost,
	}

	for _, t := range types {
		assessment := &edulab.Assessment{
			PublicID:     srv.newPublicID(2),
			Type:         t,
			ExperimentID: experiment.ID,
		}

		err = srv.DB.CreateAssessment(assessment)
		if err != nil {
			srv.renderError(w, r, err)
			return
		}
	}

	cohorts := []string{
		printer.Sprintf("Control"),
		printer.Sprintf("Intervention"),
	}

	for _, name := range cohorts {
		cohort := &edulab.Cohort{
			PublicID:     srv.newPublicID(2),
			Name:         name,
			ExperimentID: experiment.ID,
		}

		err = srv.DB.CreateCohort(cohort)
		if err != nil {
			srv.renderError(w, r, err)
			return
		}
	}

	uri := fmt.Sprintf("/experiments/%s", experiment.PublicID)
	http.Redirect(w, r, uri, http.StatusFound)
}

func (srv *Server) listExperiments(w http.ResponseWriter, r *http.Request) {
	printer, page := srv.i18n(w, r)

	experiments, err := srv.DB.FindExperiments()
	if err != nil {
		srv.renderError(w, r, err)
		return
	}
	page.Title = printer.Sprintf("Experiments")
	page.Partials = []string{"experiments"}
	page.Content = struct {
		Breadcrumbs template.HTML
		Experiments []presenter.Experiment
		Texts       interface{}
	}{
		Breadcrumbs: presenter.HomeBreadcrumbs(printer),
		Experiments: presenter.ExperimentsList(experiments, printer),
		Texts: struct {
			Title        string
			Name         string
			Participants string
			Created      string
			None         string
		}{
			Title:        printer.Sprintf("Experiments"),
			Name:         printer.Sprintf("Name"),
			Participants: printer.Sprintf("Participants"),
			Created:      printer.Sprintf("Created"),
			None:         printer.Sprintf("No available experiments"),
		},
	}
	srv.render(w, page)
}

func (srv *Server) editExperiment(w http.ResponseWriter, r *http.Request, experiment edulab.Experiment) {
	printer, page := srv.i18n(w, r)

	page.Title = printer.Sprintf("Edit Experiment: %s", experiment.Name)
	page.Partials = []string{"experiment_edit"}
	page.Content = struct {
		Breadcrumbs template.HTML
		Experiment  edulab.Experiment
		Texts       interface{}
	}{
		Breadcrumbs: presenter.ExperimentBreadcrumb(experiment, printer),
		Experiment:  experiment,
		Texts: struct {
			Edit                   string
			Name                   string
			NamePlaceholder        string
			Description            string
			DescriptionHelp        string
			DescriptionPlaceholder string
			Update                 string
		}{
			Edit:                   printer.Sprintf("Edit Experiment"),
			Name:                   printer.Sprintf("Name"),
			NamePlaceholder:        printer.Sprintf("e.g. Earth's Seasons"),
			Description:            printer.Sprintf("Description"),
			DescriptionHelp:        printer.Sprintf("Optional. Not visible to participants."),
			DescriptionPlaceholder: printer.Sprintf("e.g. This experiment will compare 2 cohorts of students. One attending a traditional lecture and the other a workshop..."),
			Update:                 printer.Sprintf("Update"),
		},
	}

	srv.render(w, page)
}

func (srv *Server) showExperiment(w http.ResponseWriter, r *http.Request, experiment edulab.Experiment) {
	printer, page := srv.i18n(w, r)

	page.Title = printer.Sprintf("Experiment: %s", experiment.Name)
	page.Partials = []string{"experiment"}
	page.Content = struct {
		Breadcrumbs template.HTML
		Experiment  edulab.Experiment
		Texts       interface{}
	}{
		Breadcrumbs: presenter.HomeBreadcrumbs(printer),
		Experiment:  experiment,
		Texts: struct {
			Experiment    string
			Settings      string
			Edit          string
			Demographics  string
			Assessments   string
			Cohorts       string
			Publish       string
			Results       string
			LearningGains string
		}{
			Experiment:    printer.Sprintf("Experiment %s", experiment.Name),
			Settings:      printer.Sprintf("Settings"),
			Edit:          printer.Sprintf("Edit Experiment"),
			Demographics:  printer.Sprintf("Demographics"),
			Assessments:   printer.Sprintf("Assessments"),
			Cohorts:       printer.Sprintf("Cohorts"),
			Publish:       printer.Sprintf("Participation Links"),
			Results:       printer.Sprintf("Results"),
			LearningGains: printer.Sprintf("Learning Gains"),
		},
	}

	srv.render(w, page)
}

func (srv *Server) updateExperiment(w http.ResponseWriter, r *http.Request, experiment edulab.Experiment) {
	err := r.ParseForm()
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	form := r.PostForm
	experiment.Name = form.Get("name")
	experiment.Description = form.Get("description")

	err = srv.DB.UpdateExperiment(experiment)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	uri := fmt.Sprintf("/experiments/%s", experiment.PublicID)
	http.Redirect(w, r, uri, http.StatusFound)
}
