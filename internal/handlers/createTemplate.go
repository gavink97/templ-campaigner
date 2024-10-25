package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/gavink97/templ-campaigner/internal/views"
	"github.com/gavink97/templ-campaigner/templates"
)

func (n *EmailDetails) CreateTemplateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		c := views.TemplateNameForm()
		if err := c.Render(r.Context(), w); err != nil {
			http.Error(w, "Failed to render content", http.StatusInternalServerError)
		}

	} else if r.Method == http.MethodPost {
		tpl := strings.ToLower(strings.TrimSuffix(r.FormValue("templatename"), ".templ"))

		headers := views.EmailHeaders{Subject: n.Headers.Subject, From: n.Headers.From, To: n.Headers.To, Component: templates.TemplateConstructor(n.Headers.To, tpl)}

		*n = EmailDetails{Headers: headers, Template: tpl, Store: n.Store}

		err := createTemplate(tpl)
		if err != nil {
			slog.Error(err.Error())
			_, err = w.Write([]byte(fmt.Sprintf("A template with that name already exists: %s", tpl)))
			if err != nil {
				slog.Info(err.Error())
			}
		}

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func createTemplate(name string) error {
	path := fmt.Sprintf("templates/%s.templ", name)

	if _, err := os.Stat(path); err == nil {
		slog.Error(err.Error())
		return err
	}

	name = MakeTitle(name)

	empty := []byte(fmt.Sprintf("package templates\n\n"+
		"templ (c *ContactDetails) %s() {\n"+"    <body>\n"+
		`        <h1 class="text-2xl font-bold">Hello { c.FName }</h1>`+
		"\n    </body>\n}", name))

	err := os.WriteFile(path, empty, 0644)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}
