package main

import (
	"io/fs"
	"net/http"
	"path/filepath"
	"text/template"
	"time"

	"github.com/justinas/nosurf"
	"od-blata-do-zlata.obradovic.dev/internal/models"
	"od-blata-do-zlata.obradovic.dev/ui"
)

type Month struct {
	Name    string
	Number  int
	Year    int
	Current bool
}

type templateData struct {
	CurrentYear     int
	CurrentMonth    int
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
	PreviousYear    int
	NextYear        int
	Year            int
	Months          []Month
	Month           int
	Incomes         []models.Income
	TotalIncome     float64
	Expenses        []models.Expense
	TotalExpense    float64
	TotalAmount     float64
	Types           []models.Type
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear:     time.Now().Year(),
		CurrentMonth:    int(time.Now().Month()),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.Format("02.01.2006")
}

func possibleNull(t *models.Type) string {
	if t == nil {
		return "---"
	}

	return t.Name
}

var functions = template.FuncMap{
	"humanDate":    humanDate,
	"possibleNull": possibleNull,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl.html",
			"html/partials/*.tmpl.html",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		// Add the template set to the map as normal...
		cache[name] = ts
	}

	return cache, nil
}
