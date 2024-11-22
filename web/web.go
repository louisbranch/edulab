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
	ID        string
	Title     string
	Header    string
	Website   string
	Layout    string
	Partials  []string
	Content   interface{}
	About     string
	FAQ       string
	ToS       string
	Languages []Language
}

type Template interface {
	Render(w io.Writer, page Page) error
}
