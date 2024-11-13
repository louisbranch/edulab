package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/louisbranch/edulab"
	"github.com/louisbranch/edulab/web/presenter"
)

func (srv *Server) publishHandler(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment) {

	log.Print("[DEBUG] web/server/publish.go: handling publish")

	printer, page := srv.i18n(w, r)

	domain := getDomainBase(r)

	cohorts, err := srv.DB.FindCohorts(experiment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	assessments, err := srv.DB.FindAssessments(experiment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	title := printer.Sprintf("Participation Links")
	page.Title = title
	page.Partials = []string{"publish"}
	page.Content = struct {
		Title       string
		Breadcrumbs template.HTML
		Domain      string
		edulab.Experiment
		Cohorts     []edulab.Cohort
		Assessments []presenter.Assessment
	}{
		Title:       title,
		Breadcrumbs: presenter.ExperimentBreadcrumb(experiment, printer),
		Domain:      domain,
		Experiment:  experiment,
		Cohorts:     cohorts,
		Assessments: presenter.NewAssessments(assessments, printer),
	}

	srv.render(w, page)
}

func getDomainBase(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	host := r.Host

	return fmt.Sprintf("%s://%s/", scheme, host)
}
