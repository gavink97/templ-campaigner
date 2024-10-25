package handlers

import (
	"net/http"

	"github.com/gavink97/templ-campaigner/internal/views"
)

func (n *EmailDetails) HideHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		c := views.ShowButton()

		if err := c.Render(r.Context(), w); err != nil {
			http.Error(w, "Failed to render content", http.StatusInternalServerError)
		}
	}
}
