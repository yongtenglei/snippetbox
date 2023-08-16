package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"rey.com/snippetbox/internal/models"
)

// Define a home handler function which writes a byte slice containing
// "Hello from Snippetbox" as the response body.
func (app application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

  data := app.newTemplateData(r)
	data.Snippets =  snippets

	app.render(w, http.StatusOK, "home.tmpl.html", data)

}

// Add a snippetView handler function.
func (app application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

  data := app.newTemplateData(r)
	data.Snippet =  snippet

	app.render(w, http.StatusOK, "view.tmpl.html", data)

}

// Add a snippetCreate handler function.
func (app application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		// set header k-v pairs
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
