package main

import (
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	type credentials struct {
		UserName string `json:"email`
		Password string `json:password`
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

	// authenticate
	app.infoLog.Println(creds.UserName, creds.Password)

	// send back a message
	payload.Error = false
	payload.Message = "signed in"

	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		app.errorLog.Println(err)
	}
}
