package server

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/louisbranch/edulab/mock"
	"github.com/louisbranch/edulab/web"
)

// TODO: make private
func ServerTest(srv *Server, req *http.Request) *httptest.ResponseRecorder {
	if srv == nil {
		srv = &Server{}
	}
	if srv.Template == nil {
		tpl := &mock.Template{}

		tpl.RenderMethod = func(w io.Writer, page web.Page) error {
			return nil
		}
		srv.Template = tpl
	}

	res := httptest.NewRecorder()
	mux := srv.NewServeMux()
	mux.ServeHTTP(res, req)

	return res
}
