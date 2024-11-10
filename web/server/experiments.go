package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/louisbranch/edulab"
	"github.com/louisbranch/edulab/web/presenter"
)

var alphanum = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func (srv *Server) experimentsHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	segments := strings.Split(strings.Trim(path, "/"), "/")

	pid := segments[0]

	if pid == "" {
		switch r.Method {
		case "GET":
			srv.listExperiments(w, r)
		case "POST":
			srv.createExperiment(w, r)
		default:
			http.NotFound(w, r)
		}
		return
	}

	if pid == "new" {
		srv.newExperimentForm(w, r)
		return
	}

	if len(segments) > 1 {
		switch segments[1] {
		case "edit":
			srv.editExperiment(w, r, pid)
			return
		case "assessments":
			srv.assessmentsForExperiment(w, r, pid)
			return
		default:
			http.NotFound(w, r)
			return
		}
	}

	srv.showExperiment(w, r, pid)
}

func (srv *Server) newExperimentForm(w http.ResponseWriter, r *http.Request) {
	printer, page := srv.i18n(w, r)

	page.Title = printer.Sprintf("New Experiment")
	page.Partials = []string{"new_experiment"}
	page.Content = struct {
		Title       string
		Name        string
		Description string
	}{
		Title:       printer.Sprintf("New Experiment"),
		Name:        printer.Sprintf("Name"),
		Description: printer.Sprintf("Description"),
	}

	srv.render(w, page)
}

func (srv *Server) createExperiment(w http.ResponseWriter, r *http.Request) {
	b := make([]rune, 6)
	for i := range b {
		b[i] = alphanum[srv.Random.Intn(len(alphanum))]
	}
	pid := fmt.Sprintf("%s-%s", string(b[:3]), string(b[3:]))

	experiment := edulab.Experiment{
		PublicID:  pid,
		CreatedAt: time.Now(),
	}

	err := r.ParseForm()
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	form := r.PostForm
	experiment.Name = form.Get("name")
	experiment.Description = form.Get("description")

	err = srv.DB.CreateExperiment(&experiment)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	uri := fmt.Sprintf("/edulab/experiments/%s", experiment.PublicID)
	http.Redirect(w, r, uri, http.StatusFound)
}

func (srv *Server) listExperiments(w http.ResponseWriter, r *http.Request) {
	printer, page := srv.i18n(w, r)

	experiments, err := srv.DB.FindExperiments()
	if err != nil {
		srv.renderError(w, r, err)
		return
	}
	page.Title = printer.Sprintf("Experiments")
	page.Partials = []string{"experiments"}
	page.Content = struct {
		Experiments []presenter.Experiment
		Name        string
		Created     string
		None        string
		Back        string
	}{
		Experiments: presenter.ExperimentsList(experiments, printer),
		Name:        printer.Sprintf("Name"),
		Created:     printer.Sprintf("Created"),
		None:        printer.Sprintf("No available experiments"),
		Back:        printer.Sprintf("Back"),
	}
	srv.render(w, page)
}

func (srv *Server) editExperiment(w http.ResponseWriter, r *http.Request, _ string) {
	srv.renderError(w, r, fmt.Errorf("not implemented"))
}

func (srv *Server) assessmentsForExperiment(w http.ResponseWriter, r *http.Request, _ string) {
	srv.renderError(w, r, fmt.Errorf("not implemented"))
}

func (srv *Server) showExperiment(w http.ResponseWriter, r *http.Request, pid string) {
	printer, page := srv.i18n(w, r)

	experiment, err := srv.DB.FindExperiment(pid)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	page.Title = printer.Sprintf("Experiment: %s", experiment.Name)
	page.Partials = []string{"experiment"}
	page.Content = struct {
		Experiment edulab.Experiment
	}{
		Experiment: experiment,
	}

	srv.render(w, page)
}
