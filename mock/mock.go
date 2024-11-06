package mock

import (
	"io"

	"github.com/louisbranch/edulab/web"
)

type Template struct {
	RenderMethod func(w io.Writer, page web.Page) error
}

func (m *Template) Render(w io.Writer, page web.Page) error {
	return m.RenderMethod(w, page)
}

type Database struct {
}
