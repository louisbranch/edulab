package server

import (
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/louisbranch/edulab"
	"github.com/louisbranch/edulab/mock"
	"github.com/louisbranch/edulab/web"
)

func serverTest(srv *Server, req *http.Request) *httptest.ResponseRecorder {
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

	if srv.DB == nil {
		srv.DB = mock.NewDB()
	}

	if srv.Random == nil {
		srv.Random = rand.New(rand.NewSource(0))
	}

	res := httptest.NewRecorder()
	mux := srv.NewServeMux()
	mux.ServeHTTP(res, req)

	return res
}

func TestGetRoutes(t *testing.T) {
	tests := []struct {
		path       string
		statusCode int
	}{
		{path: "/experiments/", statusCode: http.StatusOK},
		{path: "/experiments/new", statusCode: http.StatusOK},
		{path: "/experiments/E1", statusCode: http.StatusOK},
		{path: "/experiments/E2", statusCode: http.StatusNotFound},
		{path: "/experiments/E1/edit", statusCode: http.StatusOK},
		{path: "/experiments/E1/assessments", statusCode: http.StatusOK},
		{path: "/experiments/E1/assessments/A1", statusCode: http.StatusOK},
		{path: "/experiments/E1/assessments/A2", statusCode: http.StatusNotFound},
		{path: "/experiments/E1/assessments/A1/preview", statusCode: http.StatusOK},
		{path: "/experiments/E1/assessments/A1/questions/new", statusCode: http.StatusOK},
		{path: "/experiments/E1/assessments/A1/questions/1", statusCode: http.StatusOK},
		{path: "/experiments/E1/assessments/A1/questions/2", statusCode: http.StatusNotFound},
		{path: "/experiments/E1/demographics", statusCode: http.StatusOK},
		{path: "/experiments/E1/cohorts", statusCode: http.StatusOK},
		{path: "/experiments/E1/cohorts/C1", statusCode: http.StatusOK},
		{path: "/experiments/E1/cohorts/C2", statusCode: http.StatusNotFound},
		{path: "/experiments/E1/participate", statusCode: http.StatusOK},
		{path: "/experiments/E1/results", statusCode: http.StatusNotFound},
		{path: "/experiments/E1/results/demographics", statusCode: http.StatusOK},
		{path: "/experiments/E1/results/assessments", statusCode: http.StatusOK},
		{path: "/experiments/E1/results/gains", statusCode: http.StatusOK},
		{path: "/E1-C1-A1", statusCode: http.StatusOK},
		{path: "/E2-C1-A1", statusCode: http.StatusNotFound},
		{path: "/E1-C2-A1", statusCode: http.StatusNotFound},
		{path: "/E1-C1-A2", statusCode: http.StatusNotFound},
		{path: "/about", statusCode: http.StatusOK},
		{path: "/guide", statusCode: http.StatusOK},
		{path: "/faq", statusCode: http.StatusOK},
		{path: "/tos", statusCode: http.StatusOK},
		{path: "/", statusCode: http.StatusOK},
	}

	db := mock.NewDB()
	err := db.CreateExperiment(&edulab.Experiment{
		ID:       "1",
		PublicID: "E1",
		Name:     "Experiment 1",
	})
	if err != nil {
		t.Fatalf("failed to create experiment: %v", err)
	}

	err = db.CreateAssessment(&edulab.Assessment{
		ID:           "1",
		ExperimentID: "1",
		PublicID:     "A1",
		Type:         edulab.AssessmentTypePre,
	})
	if err != nil {
		t.Fatalf("failed to create assessment: %v", err)
	}

	err = db.CreateQuestion(&edulab.Question{
		ID:           "1",
		AssessmentID: "1",
	})
	if err != nil {
		t.Fatalf("failed to create question: %v", err)
	}

	err = db.CreateCohort(&edulab.Cohort{
		ID:           "1",
		ExperimentID: "1",
		PublicID:     "C1",
		Name:         "Control",
	})
	if err != nil {
		t.Fatalf("failed to create cohort: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.path, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			srv := &Server{
				DB: db,
			}

			res := serverTest(srv, req)
			if res.Code != tt.statusCode {
				t.Errorf("expected status %d, got %d for path %s", tt.statusCode, res.Code, tt.path)
			}
		})
	}
}
