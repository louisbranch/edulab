package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"

	"github.com/louisbranch/edulab"
	"github.com/louisbranch/edulab/result"
	"github.com/louisbranch/edulab/stats"
	"github.com/louisbranch/edulab/web/presenter"
	"gonum.org/v1/gonum/stat"
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
	case "gains":
		srv.gainsResult(w, r, experiment)
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

func (srv *Server) gainsResult(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment) {

	res, err := result.New(srv.DB, experiment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	// FIXME: This is a temporary solution to avoid a panic
	if !res.Valid() {
		err := fmt.Errorf("experiment needs at least 2 cohorts and 2 assessments")
		srv.renderError(w, r, err)
		return
	}

	if !res.Participation() {
		err := fmt.Errorf("no participation data available yet")
		srv.renderError(w, r, err)
		return
	}

	cohorts, items := res.ComparisonPairs()

	type chart struct {
		PreBase          float64 `json:"preBase"`
		PostBase         float64 `json:"postBase"`
		PreIntervention  float64 `json:"preIntervention"`
		PostIntervention float64 `json:"postIntervention"`
		Beta0            float64 `json:"beta0"`
		Beta1            float64 `json:"beta1"`
		RSquared         float64 `json:"rSquared"`
		PValue           float64 `json:"pValue"`
	}

	if r.Header.Get("Content-type") == "application/json" {
		var payload []chart

		for _, item := range items {
			comparison, err := result.NewComparison(res, item, cohorts)
			if err != nil {
				srv.renderError(w, r, err)
				return
			}

			data := comparison.ToStatsData()

			gains, interventions := stats.CalculateLearningGains(data)

			beta0, beta1, rSquared := stats.LinearRegression(gains, interventions)

			pValue := stats.ComputePValue(beta0, beta1, gains, interventions)

			if math.IsNaN(pValue) {
				pValue = 1.0
			}

			if math.IsNaN(rSquared) {
				rSquared = 0.0
			}

			var preBase, postBase, preIntervention, postIntervention []float64
			for _, d := range data {
				preBase = append(preBase, d.PreBase)
				postBase = append(postBase, d.PostBase)
				preIntervention = append(preIntervention, d.PreIntervention)
				postIntervention = append(postIntervention, d.PostIntervention)
			}

			payload = append(payload, chart{
				PreBase:          stat.Mean(preBase, nil),
				PostBase:         stat.Mean(postBase, nil),
				PreIntervention:  stat.Mean(preIntervention, nil),
				PostIntervention: stat.Mean(postIntervention, nil),
				Beta0:            beta0,
				Beta1:            beta1,
				RSquared:         rSquared,
				PValue:           pValue,
			})

		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(payload)
		if err != nil {
			srv.renderError(w, r, err)
			return
		}

		return
	}

	printer, page := srv.i18n(w, r)

	title := printer.Sprintf("Gains Results")
	page.Title = title
	page.Partials = []string{"results_gains"}
	page.Content = struct {
		Breadcrumbs template.HTML
		Experiment  edulab.Experiment
		Texts       interface{}
	}{
		Breadcrumbs: presenter.ExperimentBreadcrumb(experiment, printer),
		Experiment:  experiment,
		Texts: struct {
			Title string
			Empty string
		}{
			Title: title,
			Empty: printer.Sprintf("No data available yet"),
		},
	}

	srv.render(w, page)
}
