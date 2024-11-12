package server

import (
	"html/template"
	"log"
	"net/http"

	"github.com/louisbranch/edulab"
	"github.com/louisbranch/edulab/web/presenter"
)

func (srv *Server) cohortsHandler(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment, segments []string) {

	log.Print("[DEBUG] web/server/cohorts.go: handling cohorts")

	if len(segments) < 1 {
		if r.Method == http.MethodPost {
			srv.createCohort(w, r, experiment)
			return
		} else if r.Method == http.MethodGet {
			srv.listCohorts(w, r, experiment)
			return
		}
	}

	pid := segments[0]
	if len(segments) == 1 && r.Method == http.MethodPost {
		srv.updateCohort(w, r, experiment, pid)
		return
	}

	switch pid {
	case "new":
		srv.newCohort(w, r, experiment)
		return
	default:
		srv.showCohort(w, r, experiment, pid)
		return
	}
}

func (srv *Server) listCohorts(w http.ResponseWriter, r *http.Request, experiment edulab.Experiment) {
	printer, page := srv.i18n(w, r)

	cohorts, err := srv.DB.FindCohorts(experiment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	title := printer.Sprintf("Cohorts")
	page.Title = title
	page.Partials = []string{"cohorts"}
	page.Content = struct {
		Title       string
		Breadcrumbs template.HTML
		Experiment  edulab.Experiment
		Cohorts     []edulab.Cohort
		Texts       interface{}
	}{
		Title:       title,
		Breadcrumbs: presenter.ExperimentBreadcrumb(experiment, printer),
		Experiment:  experiment,
		Cohorts:     cohorts,
		Texts: struct {
			AddCohort string
			Name      string
			Actions   string
			Edit      string
			NoCohorts string
		}{
			AddCohort: printer.Sprintf("Add Cohort"),
			Name:      printer.Sprintf("Name"),
			Actions:   printer.Sprintf("Actions"),
			Edit:      printer.Sprintf("Edit"),
			NoCohorts: printer.Sprintf("No cohorts found"),
		},
	}

	srv.render(w, page)
}

func (srv *Server) newCohort(w http.ResponseWriter, r *http.Request, experiment edulab.Experiment) {

	printer, page := srv.i18n(w, r)

	title := printer.Sprintf("New Cohort")
	page.Title = title
	page.Partials = []string{"new_cohort"}
	page.Content = struct {
		Breadcrumbs template.HTML
		Experiment  edulab.Experiment
		Title       string
		Texts       interface{}
	}{
		Breadcrumbs: presenter.CohortBreadcrumb(experiment, printer),
		Experiment:  experiment,
		Title:       title,
		Texts: struct {
			Name                   string
			NamePlaceholder        string
			NameHelp               string
			Description            string
			DescriptionPlaceholder string
			DescriptionHelp        string
			Create                 string
		}{
			Name:                   printer.Sprintf("Name"),
			NamePlaceholder:        printer.Sprintf("e.g. Baseline"),
			NameHelp:               printer.Sprintf("Not visible to participants."),
			Description:            printer.Sprintf("Description"),
			DescriptionPlaceholder: printer.Sprintf("e.g. Cohort attending lecture-based instruction"),
			DescriptionHelp:        printer.Sprintf("Optional. Not visible to participants."),
			Create:                 printer.Sprintf("Create"),
		},
	}

	srv.render(w, page)
}

func (srv *Server) createCohort(w http.ResponseWriter, r *http.Request, experiment edulab.Experiment) {
	err := r.ParseForm()
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")

	cohort := &edulab.Cohort{
		ExperimentID: experiment.ID,
		Name:         name,
		Description:  description,
	}

	err = srv.DB.CreateCohort(cohort)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	http.Redirect(w, r, "/edulab/experiments/"+experiment.PublicID+"/cohorts", http.StatusSeeOther)
}

func (srv *Server) updateCohort(w http.ResponseWriter, r *http.Request, experiment edulab.Experiment, pid string) {
	err := r.ParseForm()
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")

	cohort := edulab.Cohort{
		Name:        name,
		Description: description,
		PublicID:    pid,
	}

	err = srv.DB.UpdateCohort(experiment.ID, cohort)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	http.Redirect(w, r, "/edulab/experiments/"+experiment.PublicID+"/cohorts/", http.StatusSeeOther)
}

func (srv *Server) showCohort(w http.ResponseWriter, r *http.Request, experiment edulab.Experiment, pid string) {

	printer, page := srv.i18n(w, r)

	cohort, err := srv.DB.FindCohort(experiment.ID, pid)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	title := printer.Sprintf("Cohort: %s", cohort.Name)
	page.Title = title
	page.Partials = []string{"cohort"}
	page.Content = struct {
		Breadcrumbs template.HTML
		Experiment  edulab.Experiment
		Cohort      edulab.Cohort
		Texts       interface{}
	}{
		Breadcrumbs: presenter.CohortBreadcrumb(experiment, printer),
		Experiment:  experiment,
		Cohort:      cohort,
		Texts: struct {
			Title           string
			Name            string
			NameHelp        string
			Description     string
			DescriptionHelp string
			Update          string
		}{
			Title:           title,
			Name:            printer.Sprintf("Name"),
			NameHelp:        printer.Sprintf("Not visible to participants."),
			Description:     printer.Sprintf("Description"),
			DescriptionHelp: printer.Sprintf("Optional. Not visible to participants."),
			Update:          printer.Sprintf("Update"),
		},
	}

	srv.render(w, page)
}
