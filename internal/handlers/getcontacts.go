package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	c "github.com/gavink97/templ-campaigner/internal/contacts"
	"github.com/gavink97/templ-campaigner/internal/views"
	"github.com/gavink97/templ-campaigner/templates"
)

func (n *EmailDetails) GetContacts(w http.ResponseWriter, r *http.Request) {
	contactList := n.Headers.To
	subject := n.Headers.Subject
	from := n.Headers.From

	headers := views.EmailHeaders{Subject: subject, From: from, To: contactList,
		Component: templates.TemplateConstructor(contactList, n.Template)}

	if r.Method == http.MethodGet {
		for _, v := range *contactList {
			component := n.Headers.ContactLabel(v.EmailAddress)
			if err := component.Render(r.Context(), w); err != nil {
				http.Error(w, "Failed to render content", http.StatusInternalServerError)
				return
			}
		}

	} else if r.Method == http.MethodPost {
		contacts := n.SearchThruContacts(r.FormValue("to"))

		// if contacts are overflown cut them off and add a number minus visible
		for _, contact := range contacts {
			for _, v := range *contactList {
				if v.EmailAddress == contact.EmailAddress {
					slog.Debug("Render contacts in list")
				}
			}

			component := headers.SearchResults(contact.EmailAddress)
			if err := component.Render(r.Context(), w); err != nil {
				http.Error(w, "Failed to render content", http.StatusInternalServerError)
				return
			}
		}

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	*n = EmailDetails{Headers: headers, Template: n.Template, Store: n.Store}
}

// needs to be refactored
func (n *EmailDetails) SearchThruContacts(email string) []c.Contact {
	var set = make(map[c.Contact]bool)

	for _, v := range *n.Headers.To {
		if set[v] {
			slog.Debug(fmt.Sprintf("%s is already in the contact list.", v.EmailAddress))
		} else {
			set[v] = true
		}
	}

	list := []c.Contact{}

	if email == "subscribed" {
		contacts, err := n.Store.GetSubscribersList()

		if err != nil {
			slog.Info("You have no subscribers in your list")
		}

		for _, v := range *contacts {
			if set[v] {
				slog.Debug(fmt.Sprintf("%s is already in the contact list.", v.EmailAddress))
			} else {
				list = append(list, v)
				set[v] = true
			}
		}
	} else if email == "" {
		list = []c.Contact{}

	} else {
		contacts, err := n.Store.SearchContacts(email)

		if err != nil {
			slog.Info(fmt.Sprintf("No names or emails match: %s", email))
		}

		for _, v := range *contacts {
			if set[v] {
				slog.Debug(fmt.Sprintf("%s is already in the contact list.", v.EmailAddress))
			} else {
				list = append(list, v)
				set[v] = true
			}
		}
	}

	return list
}
