package main

import (
	"errors"
	"net/http"
	"time"
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
