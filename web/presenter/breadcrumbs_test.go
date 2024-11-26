package presenter

import (
	"html/template"
	"testing"

	"github.com/louisbranch/edulab"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func TestHomeBreadcrumbs(t *testing.T) {
	printer := message.NewPrinter(language.English)
	expected := `<nav class="breadcrumb"><a href="/">Home</a></nav>`
	result := HomeBreadcrumbs(printer)
	if template.HTML(expected) != result {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestExperimentBreadcrumb(t *testing.T) {
	printer := message.NewPrinter(language.English)
	experiment := edulab.Experiment{Name: "Test Experiment", PublicID: "123"}
	expected := `<nav class="breadcrumb"><a href="/">Home</a> &rsaquo; <a href="/experiments/123">Test Experiment</a></nav>`
	result := ExperimentBreadcrumb(experiment, printer)
	if template.HTML(expected) != result {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestAssessmentsBreadcrumb(t *testing.T) {
	printer := message.NewPrinter(language.English)
	experiment := edulab.Experiment{Name: "Test Experiment", PublicID: "123"}
	expected := `<nav class="breadcrumb"><a href="/">Home</a> &rsaquo; <a href="/experiments/123">Test Experiment</a> &rsaquo; <a href="/experiments/123/assessments/">Assessments</a></nav>`
	result := AssessmentsBreadcrumb(experiment, printer)
	if template.HTML(expected) != result {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestAssessmentBreadcrumb(t *testing.T) {
	printer := message.NewPrinter(language.English)
	experiment := edulab.Experiment{Name: "Test Experiment", PublicID: "123"}
	assessment := edulab.Assessment{Type: edulab.AssessmentTypePre, PublicID: "456"}
	expected := `<nav class="breadcrumb"><a href="/">Home</a> &rsaquo; <a href="/experiments/123">Test Experiment</a> &rsaquo; <a href="/experiments/123/assessments">Assessments</a> &rsaquo; <a href="/experiments/123/assessments/456">Pre-Assessment</a></nav>`
	result := AssessmentBreadcrumb(experiment, assessment, printer)
	if template.HTML(expected) != result {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestCohortBreadcrumb(t *testing.T) {
	printer := message.NewPrinter(language.English)
	experiment := edulab.Experiment{Name: "Test Experiment", PublicID: "123"}
	expected := `<nav class="breadcrumb"><a href="/">Home</a> &rsaquo; <a href="/experiments/123">Test Experiment</a> &rsaquo; <a href="/experiments/123/cohorts">Cohorts</a></nav>`
	result := CohortBreadcrumb(experiment, printer)
	if template.HTML(expected) != result {
		t.Errorf("expected %s, got %s", expected, result)
	}
}
