package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gavink97/templ-campaigner/internal/views"
)

func (n *EmailDetails) NewContactHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		c := views.TemplateNewContactForm()

		if err := c.Render(r.Context(), w); err != nil {
			http.Error(w, "Failed to render content", http.StatusInternalServerError)
		}

	} else if r.Method == http.MethodPost {
		fname := r.FormValue("fname")
		lname := r.FormValue("lname")
		email := r.FormValue("emailaddress")
		subscribedForm := r.FormValue("subscribed")

		// contacts should only be unsubscribed at their request
		var subscribed bool
		var unsubscribed bool
		if subscribedForm == "" {
			subscribed = false
			unsubscribed = false
		} else {
			subscribed = true
			unsubscribed = false
		}

		_, err := n.Store.GetContact(email)
		if err != nil {
			err := n.Store.CreateContact(fname, lname, email, subscribed, unsubscribed)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)

				c := views.RegisterContactError()

				if err := c.Render(r.Context(), w); err != nil {
					http.Error(w, "Failed to render content", http.StatusInternalServerError)
				}

				return
			}
		} else {
			c := views.ContactIsRegisteredError()

			if err = c.Render(r.Context(), w); err != nil {
				http.Error(w, "error rendering template", http.StatusInternalServerError)
			}

			return
		}

		slog.Info(fmt.Sprintf("A contact is now registered: %s", email))

		return
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
