package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gavink97/templ-campaigner/internal/contacts"
	e "github.com/gavink97/templ-campaigner/internal/export"
)

// ensure contact name is going to the correct contact
func (n *EmailDetails) SendMailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := n.Headers.EmailTemplate()
		b := e.ExportTemplate(email)
		b, img := e.PrepareImages(b, "cid")

		to := formatContacts(*n.Headers.To)
		from := n.Headers.From
		subject := n.Headers.Subject

		params := e.RequestParams{To: to, From: from, Subject: subject, Body: b, Images: img}

		req := e.NewRequest(params)

		err := req.SendEmail()
		if err != nil {
			slog.Info(err.Error())
			_, err = w.Write([]byte("There was an error sending your test email"))
			if err != nil {
				slog.Info(err.Error())
			}
		} else {
			_, err = w.Write([]byte("Test email sent"))
			if err != nil {
				slog.Info(err.Error())
			}
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func formatContacts(contacts []contacts.Contact) []string {
	var list []string
	var set = make(map[string]bool)

	for _, contact := range contacts {
		email := contact.EmailAddress
		if set[email] {
			slog.Debug(fmt.Sprintf("%s already exists in fields", email))
		} else {
			list = append(list, email)
			set[email] = true
		}
	}
	return list
}
