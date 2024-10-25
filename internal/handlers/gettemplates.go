package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"unicode"

	"github.com/a-h/templ"
	"github.com/gavink97/templ-campaigner/internal/export"
	"github.com/gavink97/templ-campaigner/internal/views"
	"github.com/gavink97/templ-campaigner/templates"
	t "github.com/gavink97/templ-campaigner/templates"
)

func (n *EmailDetails) TemplateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fallback := t.TemplateConstructor(n.Headers.To, "Default")

		var tpl templ.Component
		query := r.URL.RawQuery
		if query == "" {
			tpl = fallback
		}

		files, err := os.ReadDir("templates")
		if err != nil {
			slog.Error(err.Error())
		}

		for _, file := range files {
			if strings.TrimSuffix(file.Name(), ".templ") == query {
				query = MakeTitle(query)
				tpl = t.TemplateConstructor(n.Headers.To, query)
			}
		}

		if tpl == nil {
			tpl = fallback
		}

		headers := views.EmailHeaders{Subject: n.Headers.Subject, From: n.Headers.From, To: n.Headers.To, Component: tpl}

		*n = EmailDetails{Headers: headers, Template: query, Store: n.Store}

		err = headers.EmailTemplate().Render(r.Context(), w)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
		return

	} else if r.Method == http.MethodPost {
		headers := views.EmailHeaders{Subject: n.Headers.Subject, From: n.Headers.From, To: n.Headers.To,
			Component: templates.TemplateConstructor(n.Headers.To, n.Template)}

		*n = EmailDetails{Headers: headers, Template: n.Template, Store: n.Store}

		preview := export.LivePreview(headers.PreviewTemplate())
		if err := preview.Render(r.Context(), w); err != nil {
			http.Error(w, "Failed to render content", http.StatusInternalServerError)
			return
		}
	}
}

func MakeTitle(str string) string {
	if unicode.IsLetter(rune(str[0])) {
		ex := unicode.ToUpper(rune(str[0]))
		return fmt.Sprintf("%s%s", string(ex), str[1:])
	}
	return str
}
