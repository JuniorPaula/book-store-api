package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestApplication_AllUsers(t *testing.T) {
	var mockdRows = mockedDB.NewRows([]string{"id", "email", "first_name", "last_name", "password", "active", "created_at", "updated_at", "has_token"})
	mockdRows.AddRow("1", "me@here.ca", "Jack", "Smith", "acb123", "1", time.Now(), time.Now(), "0")

	mockedDB.ExpectQuery("select \\\\* ").WillReturnRows(mockdRows)

	// create a test recoder which satissifies the requiriments for a ResponseRecorder
	rr := httptest.NewRecorder()
	// create request
	req, _ := http.NewRequest("POST", "/admin/users/all", nil)
	// call the handler
	handler := http.HandlerFunc(testApp.AllUsers)
	handler.ServeHTTP(rr, req)

	// check for expected status code
	if rr.Code != http.StatusOK {
		t.Error("AllUsers returned wrong status code of", rr.Code)
	}
}
