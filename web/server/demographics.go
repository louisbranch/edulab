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

	title := printer.Sprint("Demographics")
	page.Title = title
	page.Partials = []string{"demographics"}
	page.Content = struct {
		Title       string
		Breadcrumbs template.HTML
	}{
		Title:       title,
		Breadcrumbs: presenter.ExperimentBreadcrumb(experiment, printer),
	}

	srv.render(w, page)
}
