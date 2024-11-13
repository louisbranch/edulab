package presenter

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/louisbranch/edulab"
	"golang.org/x/text/message"
)

const breadcrumbMaxLen = 20

type Breadcrumb struct {
	URL  string
	Name string
}

func HomeBreadcrumbs(printer *message.Printer) template.HTML {
	return renderBreadcrumbs([]Breadcrumb{
		{URL: "/", Name: printer.Sprintf("Home")},
	})
}

func ExperimentBreadcrumb(e edulab.Experiment, printer *message.Printer) template.HTML {
	// truncate name for breadcrumbs
	name := e.Name
	if len(e.Name) > breadcrumbMaxLen {
		name = e.Name[:breadcrumbMaxLen] + "..."
	}

	return renderBreadcrumbs([]Breadcrumb{
		{URL: "/", Name: printer.Sprintf("Home")},
		{URL: fmt.Sprintf("/experiments/%s", e.PublicID), Name: name},
	})
}

func AssessmentsBreadcrumb(e edulab.Experiment, printer *message.Printer) template.HTML {
	return renderBreadcrumbs([]Breadcrumb{
		{URL: "/", Name: printer.Sprintf("Home")},
		{URL: fmt.Sprintf("/experiments/%s", e.PublicID), Name: e.Name},
		{URL: fmt.Sprintf("/experiments/%s/assessments/", e.PublicID), Name: printer.Sprintf("Assessments")},
	})
}

func AssessmentBreadcrumb(e edulab.Experiment, a edulab.Assessment, printer *message.Printer) template.HTML {
	ap := NewAssessment(a, printer)

	return renderBreadcrumbs([]Breadcrumb{
		{URL: "/", Name: printer.Sprintf("Home")},
		{URL: fmt.Sprintf("/experiments/%s", e.PublicID), Name: e.Name},
		{URL: fmt.Sprintf("/experiments/%s/assessments/", e.PublicID), Name: printer.Sprintf("Assessments")},
		{URL: fmt.Sprintf("/experiments/%s/assessments/%s", e.PublicID, a.PublicID), Name: ap.Type()},
	})
}

func CohortBreadcrumb(e edulab.Experiment, printer *message.Printer) template.HTML {
	return renderBreadcrumbs([]Breadcrumb{
		{URL: "/", Name: printer.Sprintf("Home")},
		{URL: fmt.Sprintf("/experiments/%s", e.PublicID), Name: e.Name},
		{URL: fmt.Sprintf("/experiments/%s/cohorts/", e.PublicID), Name: printer.Sprintf("Cohorts")},
	})
}

func renderBreadcrumbs(breadcrumbs []Breadcrumb) template.HTML {
	var sb strings.Builder
	sb.WriteString(`<nav class="breadcrumb">`)
	for i, breadcrumb := range breadcrumbs {
		if i > 0 {
			sb.WriteString(" &rsaquo; ")
		}
		sb.WriteString(fmt.Sprintf(`<a href="%s">%s</a>`, breadcrumb.URL, breadcrumb.Name))
	}
	sb.WriteString(`</nav>`)
	return template.HTML(sb.String()) // Mark the entire HTML string as safe
}
