package server

import (
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

	title := printer.Sprint("Publish Experiment")
	page.Title = title
	page.Partials = []string{"publish"}
	page.Content = struct {
		Title       string
		Breadcrumbs template.HTML
	}{
		Title:       title,
		Breadcrumbs: presenter.ExperimentBreadcrumb(experiment, printer),
	}

	srv.render(w, page)
}
