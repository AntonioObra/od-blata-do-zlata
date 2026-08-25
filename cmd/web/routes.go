package main

import (
	"net/http"

	"github.com/justinas/alice"
	"od-blata-do-zlata.obradovic.dev/ui"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	dynamic := alice.New(app.sessionManager.LoadAndSave, preventCSRF, app.authenticate)

	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /{$}", protected.ThenFunc(app.home))

	mux.Handle("GET /track/{month}", protected.ThenFunc(app.month))

	mux.Handle("GET /track/{month}/income", protected.ThenFunc(app.monthIncome))
	mux.Handle("GET /track/{month}/income/new", protected.ThenFunc(app.monthIncomeNew))
	mux.Handle("GET /track/{month}/income/{id}/edit", protected.ThenFunc(app.monthIncomeEdit))

	mux.Handle("GET /track/{month}/expense", protected.ThenFunc(app.monthExpense))
	mux.Handle("GET /track/{month}/expense/new", protected.ThenFunc(app.monthExpenseNew))
	mux.Handle("GET /track/{month}/expense/{id}/edit", protected.ThenFunc(app.monthExpenseEdit))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
