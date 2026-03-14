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

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)

	app.render(w, r, http.StatusOK, "home.html", data)
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