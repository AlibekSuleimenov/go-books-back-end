package main

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alibeksuleimenov/go-books-back-end/internal/data"
	"log"
	"os"
	"testing"
)

var testApp Application

var mockedDB sqlmock.Sqlmock

func TestMain(m *testing.M) {
	testDB, myMock, _ := sqlmock.New()
	mockedDB = myMock
	defer testDB.Close()

	testApp = Application{
		Config:      Config{},
		InfoLog:     log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog:    log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime),
		Models:      data.New(testDB),
		Environment: "development",
	}

	os.Exit(m.Run())
}
