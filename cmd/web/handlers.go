package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"od-blata-do-zlata.obradovic.dev/internal/models"
	"od-blata-do-zlata.obradovic.dev/internal/validator"
)

type incomeNewForm struct {
	Name                string  `form:"name"`
	Amount              float64 `form:"amount"`
	IncomeDate          string  `form:"incomeDate"`
	validator.Validator `form:"-"`
}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

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

	data := app.newTemplateData(r)
	data.Year = year
	data.Month = month

	app.render(w, r, http.StatusOK, "month.tmpl.html", data)
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

	data := app.newTemplateData(r)
	data.Year = year
	data.Month = month
	data.Form = incomeNewForm{
		IncomeDate: now.Format("2006-01-02"),
	}

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

		data := app.newTemplateData(r)
		data.Form = form
		data.Year = year
		data.Month = month

		app.render(w, r, http.StatusUnprocessableEntity, "month-income-new.tmpl.html", data)
		return
	}

	id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	if id == 0 {
		return
	}

	_, err = app.incomes.Insert(id, form.Name, form.Amount, incomeDate)
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

	data := app.newTemplateData(r)
	data.Year = year
	data.Month = month
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

	data := app.newTemplateData(r)
	data.Year = year
	data.Month = month
	app.render(w, r, http.StatusOK, "month-expense.tmpl.html", data)
}

func (app *application) monthExpenseEdit(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "month-expense.tmpl.html", data)
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, r, http.StatusOK, "signup.tmpl.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be an email")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be a minimum 8 chars long...")
	form.CheckField(validator.MaxBytes(form.Password, 72), "password", "This field must not be more than 72 bytes long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl.html", data)

		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		} else {
			app.serverError(w, r, err)
		}

		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, r, http.StatusOK, "login.tmpl.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be email...")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MaxBytes(form.Password, 72), "password", "This field must not be more than 72 bytes long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl.html", data)

		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		} else {
			app.serverError(w, r, err)
		}

		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
