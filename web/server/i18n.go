package server

import (
	"net/http"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	_ "github.com/louisbranch/edulab/translations"
	"github.com/louisbranch/edulab/web"
)

func (s *Server) i18n(w http.ResponseWriter, r *http.Request) (*message.Printer, web.Page) {

	var lang language.Tag

	query := r.URL.Query().Get("lang")
	cookie, err := r.Cookie("lang")
	if err != nil || query != "" {
		http.SetCookie(w, &http.Cookie{
			Name:   "lang",
			Value:  query,
			Path:   "/",
			MaxAge: 24 * 60 * 60 * 365, // 1 year
		})
	}

	if query == "" && cookie != nil {
		query = cookie.Value
	}

	switch query {
	case "pt-BR":
		lang = language.MustParse("pt-BR")
	default:
		lang = language.MustParse("en")
	}

	printer := message.NewPrinter(lang)

	languages := []web.Language{
		{Code: "en", Name: "En"},
		{Code: "pt-BR", Name: "Pt"},
	}
	for i, l := range languages {
		url := r.URL
		query := url.Query()
		query.Set("lang", l.Code)
		url.RawQuery = query.Encode()

		languages[i].URL = url.String()
	}

	page := web.Page{
		Header:    printer.Sprintf("EduLab"),
		Website:   printer.Sprintf("Astronomy Education"),
		Footer:    printer.Sprintf("About"),
		Languages: languages,
	}

	return printer, page
}
