package server

import (
	"net/http"
	"testing"

	"github.com/louisbranch/edulab/mock"
)

func TestAbout(t *testing.T) {
	req, err := http.NewRequest("GET", "/about", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	srv := &Server{}

	res := serverTest(srv, req)
	if res.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, res.Code)
	}

	page := srv.Template.(*mock.Template).PopPage()
	if page.Title != "About" {
		t.Errorf("expected title About, got %s", page.Title)
	}

	if page.Partials[0] != "about" {
		t.Errorf("expected partial about, got %s", page.Partials[0])
	}
}

func TestGuide(t *testing.T) {
	req, err := http.NewRequest("GET", "/guide", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	srv := &Server{}

	res := serverTest(srv, req)
	if res.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, res.Code)
	}

	page := srv.Template.(*mock.Template).PopPage()
	if page.Title != "Educator's Guide" {
		t.Errorf("expected title Educator's Guide, got %s", page.Title)
	}

	if page.Partials[0] != "guide" {
		t.Errorf("expected partial guide, got %s", page.Partials[0])
	}
}

func TestFAQ(t *testing.T) {
	req, err := http.NewRequest("GET", "/faq", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	srv := &Server{}

	res := serverTest(srv, req)
	if res.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, res.Code)
	}

	page := srv.Template.(*mock.Template).PopPage()
	if page.Title != "Frequently Asked Questions" {
		t.Errorf("expected title Frequently Asked Questions, got %s", page.Title)
	}

	if page.Partials[0] != "faq" {
		t.Errorf("expected partial faq, got %s", page.Partials[0])
	}
}

func TestTOS(t *testing.T) {
	req, err := http.NewRequest("GET", "/tos", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	srv := &Server{}

	res := serverTest(srv, req)
	if res.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, res.Code)
	}

	page := srv.Template.(*mock.Template).PopPage()
	if page.Title != "Terms of Service" {
		t.Errorf("expected title Terms of Service, got %s", page.Title)
	}

	if page.Partials[0] != "tos" {
		t.Errorf("expected partial tos, got %s", page.Partials[0])
	}
}
