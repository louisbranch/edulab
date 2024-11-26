package mock

import (
	"io"

	"github.com/louisbranch/edulab/web"
)

type Template struct {
	pages        []web.Page
	RenderMethod func(w io.Writer, page web.Page) error
	Err          error
}

func (t *Template) Render(w io.Writer, page web.Page) error {
	if t.RenderMethod != nil {
		return t.RenderMethod(w, page)
	}

	render := func(_ io.Writer, page web.Page) error {
		t.pages = append(t.pages, page)
		return t.Err
	}

	return render(w, page)
}

func (t *Template) PopPage() web.Page {
	if len(t.pages) == 0 {
		return web.Page{}
	}

	page := t.pages[0]
	t.pages = t.pages[1:]
	return page
}
