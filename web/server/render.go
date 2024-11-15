package server

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/louisbranch/edulab/web"
)

func (srv *Server) render(w http.ResponseWriter, page web.Page) {
	if page.Layout == "" {
		page.Layout = "layout"
	}

	err := srv.Template.Render(w, page)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
	}
}

func (srv *Server) renderError(w http.ResponseWriter, r *http.Request, err error) {
	printer, page := srv.i18n(w, r)
	if errors.Is(err, sql.ErrNoRows) {
		srv.renderNotFound(w, r)
		return
	}

	page.Title = printer.Sprintf("Internal Server Error")
	page.Content = err
	page.Partials = []string{"500"}

	w.WriteHeader(http.StatusInternalServerError)
	srv.render(w, page)
}

func (srv *Server) renderNotFound(w http.ResponseWriter, r *http.Request) {
	printer, page := srv.i18n(w, r)
	page.Title = printer.Sprintf("Page Not Found")
	page.Content = struct {
		Title string
		Home  string
	}{
		Title: printer.Sprintf("Page Not Found"),
		Home:  printer.Sprintf("Home"),
	}
	page.Partials = []string{"404"}

	w.WriteHeader(http.StatusNotFound)
	srv.render(w, page)
}
