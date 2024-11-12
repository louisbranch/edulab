package server

import (
	"html/template"
	"log"
	"net/http"

	"github.com/louisbranch/edulab"
	"github.com/louisbranch/edulab/web/presenter"
)

func (srv *Server) demographicsHandler(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment) {

	log.Print("[DEBUG] web/server/demographics.go: handling demographics")

	printer, page := srv.i18n(w, r)

	demographics, err := srv.DB.FindDemographics(experiment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	title := printer.Sprint("Demographics")
	page.Title = title
	page.Partials = []string{"demographics"}
	page.Content = struct {
		Breadcrumbs  template.HTML
		Experiment   edulab.Experiment
		Demographics []presenter.Demographic
		Texts        interface{}
	}{
		Breadcrumbs:  presenter.ExperimentBreadcrumb(experiment, printer),
		Experiment:   experiment,
		Demographics: presenter.NewDemographics(demographics, printer),
		Texts: struct {
			Title       string
			Demographic string
			Prompt      string
			Actions     string
			Add         string
			ComingSoon  string
		}{
			Title:       title,
			Demographic: printer.Sprint("Demographic"),
			Prompt:      printer.Sprint("Prompt"),
			Actions:     printer.Sprint("Actions"),
			Add:         printer.Sprint("Add Demographic"),
			ComingSoon:  printer.Sprint("Coming Soon"),
		},
	}

	srv.render(w, page)
}
