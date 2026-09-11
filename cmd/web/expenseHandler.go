package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"od-blata-do-zlata.obradovic.dev/internal/validator"
)

type expenseNewForm struct {
	Name                string  `form:"name"`
	Amount              float64 `form:"amount"`
	ExpenseDate         string  `form:"expenseDate"`
	TypeID              *int    `form:"type_id"`
	validator.Validator `form:"-"`
}

func (app *application) monthExpense(w http.ResponseWriter, r *http.Request) {
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

	expenses, err := app.expenses.List(id, year, month)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	var totalExpense float64

	for _, expense := range expenses {
		totalExpense += expense.Amount
	}

	data := app.newTemplateData(r)
	data.Year = year
	data.Month = month
	data.Expenses = expenses
	data.TotalExpense = totalExpense

	app.render(w, r, http.StatusOK, "month-expense.tmpl.html", data)
}

func (app *application) monthExpenseNew(w http.ResponseWriter, r *http.Request) {
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
	data.Form = expenseNewForm{
		ExpenseDate: now.Format("2006-01-02"),
	}
	data.Types = typesAll

	app.render(w, r, http.StatusOK, "month-expense-new.tmpl.html", data)
}

func (app *application) monthExpenseNewPost(w http.ResponseWriter, r *http.Request) {
	var form expenseNewForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	expenseDate, err := time.Parse("2006-01-02", form.ExpenseDate)
	if err != nil {
		form.AddFieldError("expenseDate", "Invalid date")
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

		app.render(w, r, http.StatusUnprocessableEntity, "month-expense-new.tmpl.html", data)
		return
	}

	id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	if id == 0 {
		return
	}

	_, err = app.expenses.Insert(id, form.Name, form.Amount, expenseDate, form.TypeID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "New expense successfully added!")

	http.Redirect(
		w,
		r,
		fmt.Sprintf(
			"/finances/%d/%d/expense",
			expenseDate.Year(),
			expenseDate.Month(),
		),
		http.StatusSeeOther,
	)
}

func (app *application) monthExpenseEdit(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "month-expense.tmpl.html", data)
}
