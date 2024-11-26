package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/louisbranch/edulab/mock"
	"github.com/louisbranch/edulab/web"
)

func TestRender(t *testing.T) {
	srv := &Server{
		Template: &mock.Template{},
	}
	page := web.Page{
		Title:    "Test",
		Content:  "Test Content",
		Partials: []string{"test"},
	}
	res := httptest.NewRecorder()
	srv.render(res, page)

	page = srv.Template.(*mock.Template).PopPage()

	if page.Title != "Test" {
		t.Errorf("expected title Test, got %s", page.Title)
	}

	template := &mock.Template{
		Err: errors.New("error"),
	}
	srv.Template = template

	res = httptest.NewRecorder()
	srv.render(res, page)

	if res.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, res.Code)
	}
}

func TestRenderError(t *testing.T) {
	srv := &Server{
		Template: &mock.Template{},
	}
	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	srv.renderError(res, req, errors.New("error"))

	page := srv.Template.(*mock.Template).PopPage()

	if page.Title != "Internal Server Error" {
		t.Errorf("expected title Internal Server Error, got %s", page.Title)
	}

	if page.Content.(error).Error() != "error" {
		t.Errorf("expected content error, got %s", page.Content)
	}

	if page.Partials[0] != "500" {
		t.Errorf("expected partial 500, got %s", page.Partials[0])
	}

	res = httptest.NewRecorder()
	srv.renderError(res, req, nil)

	if res.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, res.Code)
	}
}

func TestRenderNotFound(t *testing.T) {

	srv := &Server{
		Template: &mock.Template{},
	}
	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	srv.renderNotFound(res, req)

	page := srv.Template.(*mock.Template).PopPage()

	if page.Title != "Page Not Found" {
		t.Errorf("expected title Page Not Found, got %s", page.Title)
	}

	content := page.Content.(struct {
		Title string
		Home  string
	})
	if content.Title != "Page Not Found" {
		t.Errorf("expected content title Page Not Found, got %s", content.Title)
	}

	if content.Home != "Home" {
		t.Errorf("expected content home Home, got %s", content.Home)
	}

	if page.Partials[0] != "404" {
		t.Errorf("expected partial 404, got %s", page.Partials[0])
	}

	res = httptest.NewRecorder()
	srv.renderNotFound(res, req)

	if res.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, res.Code)
	}
}
