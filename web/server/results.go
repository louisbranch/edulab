package server

import (
	"encoding/json"
	"fmt"
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

	dp := presenter.NewDemographics(demographics, options)

	if r.Header.Get("Content-type") == "application/json" {

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

		answers, err := dp.Values(participants, participations)
		if err != nil {
			srv.renderError(w, r, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(answers)
		if err != nil {
			srv.renderError(w, r, err)
			return
		}

		return
	}

	i18n, page := srv.i18n(w, r)

	title := i18n.Sprintf("Demographics")
	page.Title = title
	page.Partials = []string{"results_demographics"}
	page.Content = struct {
		Title        string
		Experiment   edulab.Experiment
		Demographics presenter.Demographics
		Texts        interface{}
	}{
		Title:        title,
		Experiment:   experiment,
		Demographics: dp,
		Texts: struct {
			Labels       [][]string
			Options      string
			Participants string
		}{
			Labels:       dp.Labels(),
			Options:      i18n.Sprintf("Options"),
			Participants: i18n.Sprintf("Participants"),
		},
	}

	srv.render(w, page)
}
