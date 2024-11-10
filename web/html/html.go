package html

import (
	"bytes"
	"encoding/json"
	"html"
	"html/template"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/louisbranch/edulab/web"
	"github.com/yuin/goldmark"
)

type HTML struct {
	basepath string
	sync     sync.RWMutex
	cache    map[string]*template.Template
}

func New(basepath string) *HTML {
	return &HTML{
		basepath: basepath,
		cache:    make(map[string]*template.Template),
	}
}

func (h *HTML) Render(w io.Writer, page web.Page) error {
	paths := append([]string{page.Layout}, page.Partials...)

	for i, n := range paths {
		p := []string{h.basepath}
		p = append(p, strings.Split(n+".html", "/")...)
		paths[i] = filepath.Join(p...)
	}

	tpl, err := h.parse(paths...)
	if err != nil {
		return err
	}

	if page.ID == "" && len(page.Partials) > 0 {
		page.ID = page.Partials[0]
	}

	err = tpl.Execute(w, page)
	return err
}

func (h *HTML) parse(names ...string) (tpl *template.Template, err error) {
	cache := os.Getenv("APP_ENV") == "production"

	cp := make([]string, len(names))
	copy(cp, names)
	sort.Strings(cp)
	id := strings.Join(cp, ":")

	h.sync.RLock()
	tpl, ok := h.cache[id]
	h.sync.RUnlock()

	fns := template.FuncMap{
		"add":      add,
		"marshal":  marshal,
		"markdown": markdown,
	}

	if !ok {
		tpl = template.New(path.Base(names[0])).Funcs(fns)

		tpl, err = tpl.ParseFiles(names...)
		if err != nil {
			return nil, err
		}
		if cache {
			h.sync.Lock()
			h.cache[id] = tpl
			h.sync.Unlock()
		}
	}

	return tpl, nil
}

func add(a, b int) int {
	return a + b
}

func marshal(v any) string {
	s, _ := json.Marshal(v)
	return html.UnescapeString(string(s))
}

func markdown(input string) template.HTML {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(input), &buf); err != nil {
		log.Println("Error converting Markdown:", err)
		return ""
	}
	return template.HTML(buf.String()) // Safe because goldmark escapes HTML
}
