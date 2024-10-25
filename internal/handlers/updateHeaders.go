package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	c "github.com/gavink97/templ-campaigner/internal/contacts"
	"github.com/gavink97/templ-campaigner/internal/views"
	"github.com/gavink97/templ-campaigner/templates"
)

type EmailDetails struct {
	Headers  views.EmailHeaders
	Template string
	Store    c.ContactStore
}

type EmailDetailsParams struct {
	Headers  views.EmailHeaders
	Template string
	Store    c.ContactStore
}

func NewEmailDetails(params *EmailDetailsParams) *EmailDetails {
	return &EmailDetails{
		Headers:  params.Headers,
		Template: params.Template,
		Store:    params.Store,
	}
}

// create a visual indication that the details have been updated
func (n *EmailDetails) UpdateDetailsHandler(w http.ResponseWriter, r *http.Request) {
	ToHeader := n.Headers.To

	var set = make(map[c.Contact]bool)
	contactList := []c.Contact{}
	for _, v := range *ToHeader {
		if set[v] {
			slog.Debug(fmt.Sprintf("%s is already in the contact list.", v.EmailAddress))
		} else {
			contactList = append(contactList, v)
			set[v] = true
		}
	}

	var subject string
	var from string
	if r.Method == http.MethodGet {
		subject = n.Headers.Subject
		from = n.Headers.From

		headers := views.EmailHeaders{Subject: subject, From: from, To: &contactList,
			Component: templates.TemplateConstructor(&contactList, n.Template)}

		*n = EmailDetails{Headers: headers, Template: n.Template, Store: n.Store}

		// using queries is probably not the best method to solve this problem
		// these are a little buggy when being removed fix this
	} else if r.Method == http.MethodPost {
		query := r.URL.Query()
		subject = n.Headers.Subject
		from = n.Headers.From

		var value string
		if query.Has("add") {
			value = query.Get("add")
			contact, err := n.Store.GetContact(value)
			if err != nil {
				slog.Info(fmt.Sprintf("There was an error when getting this contact info: %s", err.Error()))
				return
			}

			if set[*contact] {
				slog.Debug(fmt.Sprintf("%s is already in the contact list.", contact.EmailAddress))
			} else {
				contactList = append(contactList, *contact)
				set[*contact] = true
			}

			headers := views.EmailHeaders{Subject: subject, From: from, To: &contactList,
				Component: templates.TemplateConstructor(&contactList, n.Template)}

			*n = EmailDetails{Headers: headers, Template: n.Template, Store: n.Store}

			component := headers.ContactLabel(value)
			if err := component.Render(r.Context(), w); err != nil {
				http.Error(w, "Failed to render content", http.StatusInternalServerError)
				return
			}

		} else if query.Has("remove") {
			value = query.Get("remove")
			for i, v := range contactList {
				if v.EmailAddress == value {
					contactList = append(contactList[:i], contactList[i+1:]...)
				}
			}

			headers := views.EmailHeaders{Subject: subject, From: from, To: &contactList,
				Component: templates.TemplateConstructor(&contactList, n.Template)}

			*n = EmailDetails{Headers: headers, Template: n.Template, Store: n.Store}

		} else {
			slog.Debug("Does not contain an add or remove query")
			return
		}

	} else if r.Method == http.MethodPut {
		subject = r.FormValue("subject")
		from = r.FormValue("from")

		contacts := n.SearchThruContacts(r.FormValue("to"))

		for _, contact := range contacts {
			if set[contact] {
				slog.Debug(fmt.Sprintf("%s is already in the contact list.", contact.EmailAddress))
			} else {
				contactList = append(contactList, contact)
				set[contact] = true
			}
		}

		headers := views.EmailHeaders{Subject: subject, From: from, To: &contactList,
			Component: templates.TemplateConstructor(&contactList, n.Template)}

		*n = EmailDetails{Headers: headers, Template: n.Template, Store: n.Store}

		c := headers.UpdateForm()
		if err := c.Render(r.Context(), w); err != nil {
			http.Error(w, "Failed to render content", http.StatusInternalServerError)
			return
		}

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
}
