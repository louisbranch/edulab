package server

import "net/http"

func (srv *Server) about(w http.ResponseWriter, r *http.Request) {
	printer, page := srv.i18n(w, r)
	page.Title = printer.Sprintf("About")
	page.Partials = []string{"about"}
	page.Content = struct {
		About         string
		References    string
		Context       string
		Contributions string
		Source        string
	}{
		About:         printer.Sprintf("About"),
		References:    printer.Sprintf("References"),
		Context:       printer.Sprintf(""),
		Contributions: printer.Sprintf("If you would like to contribute to the project, for example, adding more translations, get in touch:"),
		Source:        printer.Sprintf("Source Code"),
	}

	srv.render(w, page)
}

func (srv *Server) astro(w http.ResponseWriter, r *http.Request) {
	_, page := srv.i18n(w, r)
	page.Layout = "astro"

	srv.render(w, page)
}