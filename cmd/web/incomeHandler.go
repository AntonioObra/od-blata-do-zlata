package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"od-blata-do-zlata.obradovic.dev/internal/validator"
)

type incomeNewForm struct {
	Name                string  `form:"name"`
	Amount              float64 `form:"amount"`
	IncomeDate          string  `form:"incomeDate"`
	TypeID              *int    `form:"type_id"`
	validator.Validator `form:"-"`
}

func (app *application) monthIncome(w http.ResponseWriter, r *http.Request) {
	year, err := strconv.Atoi(r.PathValue("year"))
	if err != nil || year < 1 {
		http.NotFound(w, r)
		return
	}

	month, err := strconv.Atoi(r.PathValue("month"))
	if err != nil || month < 1 || month > 12 {
		http.NotFound(w, r)
		return
	}

	id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	if id == 0 {
		return
	}

	incomes, err := app.incomes.List(id, year, month)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	var totalIncome float64

	for _, income := range incomes {
		totalIncome += income.Amount
	}

	data := app.newTemplateData(r)
	data.Year = year
	data.Month = month
	data.Incomes = incomes
	data.TotalIncome = totalIncome

	app.render(w, r, http.StatusOK, "month-income.tmpl.html", data)
}

func (app *application) monthIncomeNew(w http.ResponseWriter, r *http.Request) {
	year, err := strconv.Atoi(r.PathValue("year"))
	if err != nil || year < 1 {
		http.NotFound(w, r)
		return
	}

	month, err := strconv.Atoi(r.PathValue("month"))
	if err != nil || month < 1 || month > 12 {
		http.NotFound(w, r)
		return
	}
	now := time.Now()

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
	data.Year = year
	data.Month = month
	data.Form = incomeNewForm{
		IncomeDate: now.Format("2006-01-02"),
	}
	data.Types = typesAll

	app.render(w, r, http.StatusOK, "month-income-new.tmpl.html", data)
}

func (app *application) monthIncomeNewPost(w http.ResponseWriter, r *http.Request) {
	var form incomeNewForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	incomeDate, err := time.Parse("2006-01-02", form.IncomeDate)
	if err != nil {
		form.AddFieldError("incomeDate", "Invalid date")
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Name, 100), "name", "This field cannot be more than 100 chars long...")
	form.CheckField(validator.IsPositive(form.Amount), "amount", "This field cannot be blank!")

	if !form.Valid() {
		year, err := strconv.Atoi(r.PathValue("year"))
		if err != nil || year < 1 {
			http.NotFound(w, r)
			return
		}

		month, err := strconv.Atoi(r.PathValue("month"))
		if err != nil || month < 1 || month > 12 {
			http.NotFound(w, r)
			return
		}

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
		data.Form = form
		data.Year = year
		data.Month = month
		data.Types = typesAll

		app.render(w, r, http.StatusUnprocessableEntity, "month-income-new.tmpl.html", data)
		return
	}

	id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	if id == 0 {
		return
	}

	_, err = app.incomes.Insert(id, form.Name, form.Amount, incomeDate, form.TypeID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "New income successfully added!")

	http.Redirect(
		w,
		r,
		fmt.Sprintf(
			"/track/%d/%d/income",
			incomeDate.Year(),
			incomeDate.Month(),
		),
		http.StatusSeeOther,
	)
}

func (app *application) monthIncomeEdit(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "month-income.tmpl.html", data)
}
