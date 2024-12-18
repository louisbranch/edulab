package server

import (
	"math/rand"
	"net/http"
	"strings"

	"github.com/louisbranch/edulab"
	"github.com/louisbranch/edulab/web"
)

type Server struct {
	DB       edulab.Database
	Template web.Template
	Assets   http.Handler
	Random   *rand.Rand
}

func (srv *Server) NewServeMux() *http.ServeMux {

	mux := http.NewServeMux()

	// Static assets
	mux.Handle("/assets/", http.StripPrefix("/assets/", srv.Assets))

	// Dynamic routes
	mux.HandleFunc("/experiments/", srv.experimentsHandler)
	mux.HandleFunc("/demographics", srv.participateDemographics)
	mux.HandleFunc("/assessments", srv.participateAssessments)

	mux.HandleFunc("/about", srv.about)
	mux.HandleFunc("/guide", srv.guide)
	mux.HandleFunc("/faq", srv.faq)
	mux.HandleFunc("/tos", srv.tos)
	mux.HandleFunc("/", srv.index)

	return mux
}

func (srv *Server) index(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path
	segments := strings.Split(strings.Trim(path, "/"), "/")

	if len(segments) > 1 {
		srv.renderNotFound(w, r)
		return
	}

	pid := segments[0]

	if pid != "" {
		srv.participationsHandler(w, r, pid)
		return
	}

	printer, page := srv.i18n(w, r)
	page.Title = printer.Sprintf("EduLab - Empowering Educators")
	page.Partials = []string{"index"}
	page.Content = struct {
		Tagline             string
		Introduction        string
		Paper               string
		NewExperiment       string
		EducatorsGuide      string
		PreviousExperiments string
	}{
		Tagline: printer.Sprintf("Empowering Educators Through Evidence-Based Insights"),
		Introduction: printer.Sprintf(`EduLab brings **data-driven** experimentation into the classroom, empowering you to evaluate and refine teaching methods across distinct **cohorts**.

By running controlled pre- and post-assessments, you gain **evidence-based insights** into how different teaching approaches impact learning outcomes.

Compare cohorts, **measure learning gains**, and adapt strategies to elevate student engagement—all supported by real-time educational data.`),
		Paper:               printer.Sprintf("Read our draft paper:"),
		NewExperiment:       printer.Sprintf("New Experiment"),
		EducatorsGuide:      printer.Sprintf("Educator's Guide"),
		PreviousExperiments: printer.Sprintf("Previous Experiments"),
	}

	srv.render(w, page)
}

var alphanum = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func (srv *Server) newPublicID(lens ...int) string {
	sum := 0
	for _, l := range lens {
		sum += l
	}

	b := make([]rune, sum)
	for i := range b {
		b[i] = alphanum[srv.Random.Intn(len(alphanum))]
	}

	pid := ""
	for i, l := range lens {
		pid += string(b[:l])
		if i < len(lens)-1 {
			pid += "-"
		}
		b = b[l:]
	}

	return pid
}
