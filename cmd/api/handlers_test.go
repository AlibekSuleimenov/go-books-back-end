package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestApplication_AllUsers(t *testing.T) {
	// create mock rows, and add one row
	mockedRows := mockedDB.NewRows([]string{"id", "email", "first_name", "last_name", "password", "active", "created_at", "updated-at", "has_token"})
	mockedRows.AddRow("1", "me@here.com", "Jack", "Smith", "abc123", "1", time.Now(), time.Now(), "0")

	// tell mock what queries we expect
	mockedDB.ExpectQuery("select \\\\* ").WillReturnRows(mockedRows)

	// create a test recorder which satisfies the requirements of a ResponseRecorder
	rr := httptest.NewRecorder()
	// create a request
	request, _ := http.NewRequest("POST", "admin/users", nil)
	// call the handler
	handler := http.HandlerFunc(testApp.AllUsers)

	handler.ServeHTTP(rr, request)

	// check for expected status code
	if rr.Code != http.StatusOK {
		t.Error("AllUsers returned wrong status code of", rr.Code)
	}
}
