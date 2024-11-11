package server

import (
	"html/template"
	"log"
	"net/http"

	"github.com/louisbranch/edulab"
	"github.com/louisbranch/edulab/web/presenter"
)

func (srv *Server) cohortsHandler(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment) {

	log.Print("[DEBUG] web/server/cohorts.go: handling cohorts")

	printer, page := srv.i18n(w, r)

	title := printer.Sprint("Cohorts")
	page.Title = title
	page.Partials = []string{"cohorts"}
	page.Content = struct {
		Title       string
		Breadcrumbs template.HTML
	}{
		Title:       title,
		Breadcrumbs: presenter.ExperimentBreadcrumb(experiment, printer),
	}

	srv.render(w, page)
}
