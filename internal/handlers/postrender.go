package handlers

import (
	"fmt"
	"net/http"

	"github.com/gavink97/templ-campaigner/internal/export"
	"github.com/gavink97/templ-campaigner/internal/views"
)

func (n *EmailDetails) RenderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := n.Headers.EmailTemplate()
		s := export.ExportTemplate(email)
		fmt.Println(s)

		// w.Write([]byte("Template rendered"))
		c := views.RenderButton()
		if err := c.Render(r.Context(), w); err != nil {
			http.Error(w, "Failed to render content", http.StatusInternalServerError)
		}

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
