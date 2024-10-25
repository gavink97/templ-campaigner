package handlers

import (
	"net/http"
)

func (n *EmailDetails) ShowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		c := n.Headers.UpdateForm()

		if err := c.Render(r.Context(), w); err != nil {
			http.Error(w, "Failed to render content", http.StatusInternalServerError)
		}
	}
}
