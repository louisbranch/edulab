package server

import (
	"fmt"
	"math/rand"
	"net/http"

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

	mux.Handle("/edulab/assets/", http.StripPrefix("/edulab/assets/", srv.Assets))

	mux.HandleFunc("/edulab/about/", srv.about)

	mux.HandleFunc("/edulab/", srv.index)
	mux.HandleFunc("/", srv.astro)

	return mux
}

func (srv *Server) index(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/edulab/"):]

	if name != "" {
		fmt.Printf("404: %s\n", name)
		return
	}

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	printer, page := srv.i18n(w, r)
	page.Title = printer.Sprintf("EduLab")
	page.Partials = []string{"index"}
	page.Content = struct {
	}{}

	srv.render(w, page)
}
