package server

import (
	"html/template"
	"net/http"

	"github.com/louisbranch/edulab/web/presenter"
)

func (srv *Server) about(w http.ResponseWriter, r *http.Request) {
	printer, page := srv.i18n(w, r)
	title := printer.Sprintf("About")
	page.Title = title
	page.Partials = []string{"about"}
	page.Content = struct {
		Breadcrumbs   template.HTML
		Title         string
		References    string
		Context       string
		Contributions string
		Source        string
	}{
		Breadcrumbs:   presenter.HomeBreadcrumbs(printer),
		Title:         title,
		References:    printer.Sprintf("References"),
		Context:       printer.Sprintf(""),
		Contributions: printer.Sprintf("If you would like to contribute to the project, for example, adding more translations, get in touch:"),
		Source:        printer.Sprintf("Source Code"),
	}

	srv.render(w, page)
}

func (srv *Server) guide(w http.ResponseWriter, r *http.Request) {
	printer, page := srv.i18n(w, r)

	title := printer.Sprintf("Educator's Guide")
	page.Title = title
	page.Partials = []string{"guide"}
	page.Content = struct {
		Breadcrumbs template.HTML
		Title       string
	}{
		Breadcrumbs: presenter.HomeBreadcrumbs(printer),
		Title:       title,
	}

	srv.render(w, page)
}

func (srv *Server) faq(w http.ResponseWriter, r *http.Request) {
	printer, page := srv.i18n(w, r)

	title := printer.Sprintf("Frequently Asked Questions")
	page.Title = title
	page.Partials = []string{"faq"}
	page.Content = struct {
		Breadcrumbs template.HTML
		Title       string
	}{
		Breadcrumbs: presenter.HomeBreadcrumbs(printer),
		Title:       title,
	}

	srv.render(w, page)
}
