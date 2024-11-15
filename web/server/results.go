package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/louisbranch/edulab"
	"github.com/louisbranch/edulab/web/presenter"
)

func (srv *Server) resultsHandler(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment, segments []string) {

	log.Print("[DEBUG] Routing results")

	if len(segments) < 1 {
		srv.listResults(w, r, experiment)
		return
	}

	switch segments[0] {
	case "demographics":
		srv.demographicsResult(w, r, experiment)
		return
	default:
		srv.renderNotFound(w, r)
		return
	}
}

func (s *Server) listResults(w http.ResponseWriter, _ *http.Request,
	experiment edulab.Experiment) {

	log.Print("[DEBUG] Listing results")

	fmt.Fprintf(w, "Listing results for experiment %s", experiment.ID)
}

func (srv *Server) demographicsResult(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment) {

	cohorts, err := srv.DB.FindCohorts(experiment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	demographics, err := srv.DB.FindDemographics(experiment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	options, err := srv.DB.FindDemographicOptions(experiment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	participants, err := srv.DB.FindParticipants(experiment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	participations, err := srv.DB.FindParticipations(experiment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	dr, err := presenter.NewDemographicsResult(demographics, options, cohorts,
		participants, participations)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	if r.Header.Get("Content-type") == "application/json" {

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(dr.Data)
		if err != nil {
			srv.renderError(w, r, err)
			return
		}

		return
	}

	printer, page := srv.i18n(w, r)

	title := printer.Sprintf("Demographics Results")
	page.Title = title
	page.Partials = []string{"results_demographics"}
	page.Content = struct {
		Breadcrumbs template.HTML
		Experiment  edulab.Experiment
		Results     presenter.DemographicsResult
		Texts       interface{}
	}{
		Breadcrumbs: presenter.ExperimentBreadcrumb(experiment, printer),
		Experiment:  experiment,
		Results:     dr,
		Texts: struct {
			Title        string
			Options      string
			Participants string
			Empty        string
		}{
			Title:        title,
			Options:      printer.Sprintf("Options"),
			Participants: printer.Sprintf("Participants"),
			Empty:        printer.Sprintf("No data available yet"),
		},
	}

	srv.render(w, page)
}
