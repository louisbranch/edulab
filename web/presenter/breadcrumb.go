package presenter

import (
	"fmt"
	"html/template"
	"math"
	"strings"

	"github.com/louisbranch/edulab"
	"golang.org/x/text/message"
)

type Breadcrumb struct {
	URL  string
	Name string
}

func HomeBreadcrumbs(printer *message.Printer) template.HTML {
	return RenderBreadcrumbs([]Breadcrumb{
		{URL: "/edulab/", Name: printer.Sprintf("Home")},
	})
}

func ExperimentBreadcrumb(e edulab.Experiment, printer *message.Printer) template.HTML {
	// truncate name for breadcrumbs
	name := e.Name[:int(math.Min(20, float64(len(e.Name))))] + "..."

	return RenderBreadcrumbs([]Breadcrumb{
		{URL: "/edulab/", Name: printer.Sprintf("Home")},
		{URL: fmt.Sprintf("/edulab/experiments/%s", e.PublicID), Name: name},
	})
}

func AssessmentBreadcrumb(e edulab.Experiment, a edulab.Assessment, printer *message.Printer) template.HTML {
	// truncate name for breadcrumbs
	name := a.Name[:int(math.Min(20, float64(len(a.Name))))] + "..."

	return RenderBreadcrumbs([]Breadcrumb{
		{URL: "/edulab/", Name: printer.Sprintf("Home")},
		{URL: fmt.Sprintf("/edulab/experiments/%s", e.PublicID), Name: e.Name},
		{URL: fmt.Sprintf("/edulab/experiments/%s/assessments/%s", e.PublicID, a.PublicID), Name: name},
	})
}

func CohortBreadcrumb(e edulab.Experiment, printer *message.Printer) template.HTML {
	return RenderBreadcrumbs([]Breadcrumb{
		{URL: "/edulab/", Name: printer.Sprintf("Home")},
		{URL: fmt.Sprintf("/edulab/experiments/%s", e.PublicID), Name: e.Name},
		{URL: fmt.Sprintf("/edulab/experiments/%s/cohorts/", e.PublicID), Name: printer.Sprint("Cohorts")},
	})
}

func RenderBreadcrumbs(breadcrumbs []Breadcrumb) template.HTML {
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
