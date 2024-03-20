package main

import (
	"books_api/internal/data"
	"log"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

var testApp application
var mockedDB sqlmock.Sqlmock

func TestMain(m *testing.M) {
	testDB, myMock, _ := sqlmock.New()
	mockedDB = myMock
	defer testDB.Close()

	testApp = application{
		config:   config{},
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		errorLog: log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime),
		models:   data.New(testDB),
	}

	os.Exit(m.Run())
}
