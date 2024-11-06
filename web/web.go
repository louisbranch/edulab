package web

import (
	"io"
)

type Language struct {
	Name string
	Code string
	URL  string
}

type Page struct {
	Title     string
	Header    string
	Website   string
	Layout    string
	Partials  []string
	Content   interface{}
	Footer    string
	Languages []Language
}

type Template interface {
	Render(w io.Writer, page Page) error
}
