package routes

import (
	"net/http"

	"github.com/gavink97/templ-campaigner/internal/contacts"
	"github.com/gavink97/templ-campaigner/internal/contacts/store"
	h "github.com/gavink97/templ-campaigner/internal/handlers"
	"github.com/gavink97/templ-campaigner/internal/views"
	t "github.com/gavink97/templ-campaigner/templates"
)

func newRouter() http.Handler {
	mux := http.NewServeMux()

	db := store.MustOpen("contacts.db")

	contactStore := store.NewContactStore(store.NewContactStoreParams{
		DB: db,
	})

	contactlist := []contacts.Contact{}

	headers := views.EmailHeaders{
		Subject:   "Subject",
		From:      "noreply@gav.ink",
		To:        &contactlist,
		Component: t.TemplateConstructor(&contactlist, "Default"),
	}

	emailParams := h.EmailDetails{
		Headers:  headers,
		Template: "Default",
		Store:    contactStore,
	}

	publicFiles := http.FileServer(http.Dir("./public"))
	mux.Handle("/public/", http.StripPrefix("/public/", publicFiles))

	publicAssets := http.FileServer(http.Dir("./assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", publicAssets))

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			return
		}
		emailParams.HomeHandler(w, r)
	}))

	mux.HandleFunc("POST /render", http.HandlerFunc(emailParams.RenderHandler))

	mux.Handle("POST /send", http.HandlerFunc(emailParams.SendMailHandler))

	mux.Handle("GET /create", http.HandlerFunc(emailParams.CreateTemplateHandler))
	mux.Handle("POST /create", http.HandlerFunc(emailParams.CreateTemplateHandler))

	mux.Handle("GET /update", http.HandlerFunc(emailParams.UpdateDetailsHandler))
	mux.Handle("POST /update", http.HandlerFunc(emailParams.UpdateDetailsHandler))
	mux.Handle("PUT /update", http.HandlerFunc(emailParams.UpdateDetailsHandler))

	mux.Handle("GET /template", http.HandlerFunc(emailParams.TemplateHandler))
	mux.Handle("POST /template", http.HandlerFunc(emailParams.TemplateHandler))

	mux.Handle("GET /newcontact", http.HandlerFunc(emailParams.NewContactHandler))
	mux.Handle("POST /newcontact", http.HandlerFunc(emailParams.NewContactHandler))

	mux.Handle("GET /contacts", http.HandlerFunc(emailParams.GetContacts))
	mux.Handle("POST /contacts", http.HandlerFunc(emailParams.GetContacts))

	mux.Handle("POST /hide", http.HandlerFunc(emailParams.HideHandler))
	mux.Handle("POST /show", http.HandlerFunc(emailParams.ShowHandler))

	return mux
}
