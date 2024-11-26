package server

import (
	"net/http/httptest"
	"testing"

	"github.com/louisbranch/edulab/mock"
)

func TestI18n(t *testing.T) {
	srv := &Server{
		Template: &mock.Template{},
	}
	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	printer, page := srv.i18n(res, req)

	if printer == nil {
		t.Error("expected printer, got nil")
	}

	if page.About != "About" {
		t.Errorf("expected About, got %s", page.About)
	}

	res = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/?lang=pt-BR", nil)
	_, page = srv.i18n(res, req)

	if page.About != "Sobre" {
		t.Errorf("expected Sobre, got %s", page.About)
	}

	cookies := res.Result().Cookies()

	if len(cookies) != 1 {
		t.Errorf("expected 1 cookie, got %d", len(cookies))
	}

	if cookies[0].Name != "lang" {
		t.Errorf("expected lang, got %s", cookies[0].Name)
	}

	if cookies[0].Value != "pt-BR" {
		t.Errorf("expected pt-BR, got %s", cookies[0].Value)
	}
}
