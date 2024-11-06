package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/louisbranch/edulab"
	"github.com/louisbranch/edulab/web/presenter"
)

var alphanum = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func (srv *Server) newExperiment() edulab.Experiment {
	b := make([]rune, 6)
	for i := range b {
		b[i] = alphanum[srv.Random.Intn(len(alphanum))]
	}
	name := fmt.Sprintf("%s-%s", string(b[:3]), string(b[3:]))

	experiment := edulab.Experiment{
		Name:      name,
		CreatedAt: time.Now(),
	}

	return experiment
}

func (srv *Server) experiments(w http.ResponseWriter, r *http.Request) {
	printer, page := srv.i18n(w, r)

	switch r.Method {
	case "GET":

		name := r.URL.Path[len("/edulab/experiments/"):]
		if name == "" {
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
			return
		}

		experiment, err := srv.DB.FindExperiment(name)
		if err != nil {
			srv.renderError(w, r, err)
			return
		}

		page.Title = printer.Sprintf("Experiment %s", name)
		page.Partials = []string{"experiment"}
		page.Content = struct {
			Experiment edulab.Experiment
		}{
			Experiment: experiment,
		}

		srv.render(w, page)

	case "POST":
		experiment := srv.newExperiment()

		err := srv.DB.CreateExperiment(&experiment)
		if err != nil {
			srv.renderError(w, r, err)
			return
		}

		uri := fmt.Sprintf("/edulab/experiments/%s", experiment.Name)

		http.Redirect(w, r, uri, http.StatusFound)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
