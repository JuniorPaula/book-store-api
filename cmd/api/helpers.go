package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1048576 // 1 mega byte
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single json value")
	}

	return nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var customErr error

	switch {
	case strings.Contains(err.Error(), "SQLSTATE 23505"):
		customErr = errors.New("duplicate value violates unique constraint")
		statusCode = http.StatusForbidden
	case strings.Contains(err.Error(), "SQLSTATE 22001"):
		customErr = errors.New("the value you are trying to insert is too large")
		statusCode = http.StatusForbidden
	case strings.Contains(err.Error(), "SQLSTATE 23503"):
		customErr = errors.New("foreign key violation")
		statusCode = http.StatusForbidden
	default:
		customErr = err
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = customErr.Error()

	app.writeJSON(w, statusCode, payload)

	return nil
}
