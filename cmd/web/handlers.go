package main

import "net/http"

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// snippets, err := app.snippets.Latest()
	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	return
	// }

	data := app.newTemplateData(r)
	// data.Snippets = snippets

	app.render(w, r, http.StatusOK, "home.tmpl.html", data)
}
