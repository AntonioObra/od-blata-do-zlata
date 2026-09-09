package main

import (
	"net/http"

	"od-blata-do-zlata.obradovic.dev/internal/validator"
)

type typesNewForm struct {
	Name                string `form:"name"`
	validator.Validator `form:"-"`
}

func (app *application) typesAll(w http.ResponseWriter, r *http.Request) {
	id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	if id == 0 {
		return
	}

	typesAll, err := app.types.List(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Types = typesAll

	app.render(w, r, http.StatusOK, "types-all.tmpl.html", data)
}

func (app *application) typesNew(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = typesNewForm{}
	app.render(w, r, http.StatusOK, "types-new.tmpl.html", data)
}

func (app *application) typesNewPost(w http.ResponseWriter, r *http.Request) {
	var form typesNewForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Name, 100), "name", "This field cannot be more than 100 chars long...")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form

		app.render(w, r, http.StatusUnprocessableEntity, "types-new.tmpl.html", data)
		return
	}

	id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	if id == 0 {
		return
	}

	_, err = app.types.Insert(id, form.Name)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "New type successfully added!")

	http.Redirect(
		w,
		r,
		"/track/types",
		http.StatusSeeOther,
	)
}
