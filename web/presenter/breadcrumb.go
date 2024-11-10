package presenter

import (
	"fmt"
	"html/template"
	"strings"
)

type Breadcrumb struct {
	URL  string
	Name string
}

func BreadcrumbsHome() template.HTML {
	return RenderBreadcrumbs([]Breadcrumb{
		{URL: "/edulab/", Name: "Home"},
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
