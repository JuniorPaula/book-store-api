package main

import (
	"books_api/internal/data"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type jsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type envelope map[string]interface{}

func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	type credentials struct {
		UserName string `json:"email"`
		Password string `json:"password"`
	}

	var creds credentials
	var payload jsonResponse

	err := app.readJSON(w, r, &creds)
	if err != nil {
		app.errorLog.Println(err)
		payload.Error = true
		payload.Message = "internal server error"
		_ = app.writeJSON(w, http.StatusBadRequest, payload)
	}

	user, err := app.models.User.GetByEmail(creds.UserName)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	validPassword, err := user.PasswordMatches(creds.Password)
	if err != nil || !validPassword {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	token, err := app.models.Token.GenerateToken(user.ID, 24*time.Hour)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.models.Token.Insert(*token, *user)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload = jsonResponse{
		Error:   false,
		Message: "logged in",
		Data:    envelope{"token": token, "user": user},
	}

	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) Logout(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Token string `json:"token"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, errors.New("invalid json"), http.StatusUnprocessableEntity)
		return
	}

	err = app.models.Token.DeleteByToken(requestPayload.Token)
	if err != nil {
		app.errorJSON(w, errors.New("internal server error"), http.StatusInternalServerError)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "logout successfully",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

// AllUsers get the users from the database and return them
func (app *application) AllUsers(w http.ResponseWriter, r *http.Request) {
	var users data.User
	all, err := users.GetAll()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "success",
		Data:    envelope{"users": all},
	}

	app.writeJSON(w, http.StatusOK, payload)
}

// EditUser its a handler to edit or save a new user
// if userID great than 0, edit an user, then save a new user.
func (app *application) EditUser(w http.ResponseWriter, r *http.Request) {
	var user data.User
	err := app.readJSON(w, r, &user)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if user.ID == 0 {
		// add user
		if _, err := app.models.User.Insert(user); err != nil {
			app.errorJSON(w, err)
			return
		}
	} else {
		// edit user
		u, err := app.models.User.GetById(user.ID)
		if err != nil {
			app.errorJSON(w, err)
			return
		}

		u.Email = user.Email
		u.FirstName = user.FirstName
		u.LastName = user.LastName

		if err := u.Update(); err != nil {
			app.errorJSON(w, err)
			return
		}

		if user.Password != "" {
			err := u.ResetPassword(user.Password)
			if err != nil {
				app.errorJSON(w, err)
				return
			}
		}
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Change saved",
	}

	_ = app.writeJSON(w, http.StatusCreated, payload)
}

// GetUser get a specific user from the database
func (app *application) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	user, err := app.models.User.GetById(userID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, user)
}
