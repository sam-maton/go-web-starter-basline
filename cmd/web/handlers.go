package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/sam-maton/go-web-starter-baseline/internal/models"
	"github.com/sam-maton/go-web-starter-baseline/internal/validator"
)

type userCreateForm struct {
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

type todoCreateForm struct {
	Title               string `form:"title"`
	validator.Validator `form:"-"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)

	todos, err := app.todos.InProgress()

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data.Todos = todos

	app.render(w, r, http.StatusOK, "home.html", data)
}

func (app *application) todoCreatePost(w http.ResponseWriter, r *http.Request) {

	var form todoCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = app.todos.Insert(form.Title)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), FLASH_KEY, "Todo was successfully created!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) todoDeletePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.notFound(w)
	}

	fmt.Println(id)

	app.sessionManager.Put(r.Context(), FLASH_KEY, "Todo was successfully deleted!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userCreateForm{}
	app.render(w, r, http.StatusOK, "signup.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", NOT_BLANK_ERROR)
	form.CheckField(validator.NotBlank(form.Email), "email", NOT_BLANK_ERROR)
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Must be a valid email")
	form.CheckField(validator.NotBlank(form.Password), "password", NOT_BLANK_ERROR)
	form.CheckField(validator.MinChars(form.Password, 8), "password", "Password must be at least 8 characters")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, SIGNUP_PAGE, data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, SIGNUP_PAGE, data)
		} else {
			app.serverError(w, r, err)
		}

		return
	}

	app.sessionManager.Put(r.Context(), FLASH_KEY, "Your account has been created successfully! Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, r, http.StatusOK, "login.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, LOGIN_PAGE, data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.html", data)
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

	app.sessionManager.Put(r.Context(), AUTH_USER_KEY, id)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if app.sessionManager.Exists(r.Context(), AUTH_USER_KEY) {
		app.sessionManager.Remove(r.Context(), AUTH_USER_KEY)
		app.sessionManager.Put(r.Context(), FLASH_KEY, "You have been logged out successfully")
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
