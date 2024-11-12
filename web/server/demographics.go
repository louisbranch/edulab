package server

import (
	"html/template"
	"log"
	"net/http"

	"github.com/louisbranch/edulab"
	"github.com/louisbranch/edulab/web/presenter"
)

func (srv *Server) demographicsHandler(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment, segments []string) {

	log.Print("[DEBUG] web/server/demographics.go: handling demographics")

	if len(segments) < 1 {
		srv.listDemographics(w, r, experiment)
		return
	}

	switch segments[0] {
	case "preview":
		srv.previewDemographics(w, r, experiment)
		return
	default:
		srv.renderNotFound(w, r)
		return
	}
}

func (srv *Server) listDemographics(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment) {

	printer, page := srv.i18n(w, r)

	demographics, err := srv.DB.FindDemographics(experiment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	title := printer.Sprintf("Demographics")
	page.Title = title
	page.Partials = []string{"demographics"}
	page.Content = struct {
		Breadcrumbs  template.HTML
		Experiment   edulab.Experiment
		Demographics []edulab.Demographic
		Texts        interface{}
	}{
		Breadcrumbs:  presenter.ExperimentBreadcrumb(experiment, printer),
		Experiment:   experiment,
		Demographics: demographics,
		Texts: struct {
			Title       string
			Demographic string
			Text        string
			Actions     string
			Add         string
			ComingSoon  string
			Preview     string
		}{
			Title:       title,
			Demographic: printer.Sprintf("Demographic"),
			Text:        printer.Sprintf("Text"),
			Actions:     printer.Sprintf("Actions"),
			Add:         printer.Sprintf("Add Demographic"),
			ComingSoon:  printer.Sprintf("Coming Soon"),
			Preview:     printer.Sprintf("Preview Questions"),
		},
	}

	srv.render(w, page)
}

func (srv *Server) previewDemographics(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment) {

	printer, page := srv.i18n(w, r)

	demographics, err := srv.DB.FindDemographics(experiment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	dp := make(map[string]presenter.Demographic)
	for _, d := range demographics {
		dp[d.ID] = presenter.Demographic{
			Demographic: d,
		}
	}

	options, err := srv.DB.FindDemographicOptions(experiment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	for _, o := range options {
		d, ok := dp[o.DemographicID]
		if !ok {
			log.Printf("[ERROR] web/server/assessments.go: demographic not found: %s", o.DemographicID)
			continue
		}
		d.Options = append(d.Options, o)
		dp[o.DemographicID] = d
	}

	page.Title = printer.Sprintf("Demographics")
	page.Partials = []string{"demographics_preview"}
	page.Content = struct {
		Breadcrumbs  template.HTML
		Experiment   edulab.Experiment
		Demographics []presenter.Demographic
		Texts        interface{}
	}{
		Breadcrumbs:  presenter.ExperimentBreadcrumb(experiment, printer),
		Experiment:   experiment,
		Demographics: presenter.SortDemographics(demographics, dp),
		Texts: struct {
			Title  string
			Submit string
			Back   string
		}{
			Title:  printer.Sprintf("Demographics"),
			Submit: printer.Sprintf("Submit"),
			Back:   printer.Sprintf("Back"),
		},
	}

	srv.render(w, page)
}