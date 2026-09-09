package main

import (
	"net/http"
	"strconv"
	"time"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	now := time.Now()
	year := now.Year()
	previousYear := year - 1
	nextYear := year + 1

	if r.URL.Query().Has("year") {
		yearParam := r.URL.Query().Get("year")

		parsedYear, err := strconv.Atoi(yearParam)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		year = parsedYear
		previousYear = year - 1
		nextYear = year + 1
	}

	var months []Month

	for month := 1; month <= 12; month++ {
		months = append(months, Month{
			Name:    time.Month(month).String(),
			Number:  month,
			Year:    year,
			Current: month == int(now.Month()),
		})
	}

	data.Months = months
	data.Year = year
	data.PreviousYear = previousYear
	data.NextYear = nextYear

	app.render(w, r, http.StatusOK, "home.tmpl.html", data)
}
