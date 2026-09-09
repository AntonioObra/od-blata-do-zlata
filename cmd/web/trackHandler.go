package main

import (
	"net/http"
	"strconv"
)

func (app *application) month(w http.ResponseWriter, r *http.Request) {
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

	totalIncome, err := app.incomes.GetTotal(id, year, month)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	totalExpense, err := app.expenses.GetTotal(id, year, month)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Year = year
	data.Month = month
	data.TotalIncome = totalIncome
	data.TotalExpense = totalExpense
	data.TotalAmount = totalIncome - totalExpense

	app.render(w, r, http.StatusOK, "month.tmpl.html", data)
}
