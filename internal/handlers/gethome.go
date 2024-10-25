package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gavink97/templ-campaigner/internal/export"
	"github.com/gavink97/templ-campaigner/internal/views"
)

// https://medium.com/@dhanushgopinath/sending-html-emails-using-templates-in-golang-9e953ca32f3d
func (n *EmailDetails) HomeHandler(w http.ResponseWriter, r *http.Request) {

	preview := export.LivePreview(n.Headers.PreviewTemplate())
	header := views.EmailHeaders{To: n.Headers.To, From: n.Headers.From, Subject: n.Headers.Subject, Component: preview}

	err := header.Index().Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		slog.Error(err.Error())
		return
	}
}
