package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_readJSON(t *testing.T) {
	// create sample JSON
	sampleJSON := map[string]interface{}{
		"foo": "bar",
	}

	body, _ := json.Marshal(sampleJSON)

	// declare a variable that we can read into
	var decodeJSON struct {
		Foo string `json:"foo"`
	}

	// create a request
	request, err := http.NewRequest("POST", "/", bytes.NewReader(body))
	if err != nil {
		t.Log(err)
	}

	// create a test ResponseRecorder
	rr := httptest.NewRecorder()
	defer request.Body.Close()

	err = testApp.readJSON(rr, request, &decodeJSON)
	if err != nil {
		t.Error("failed to decode json", err)
	}
}

func Test_writeJSON(t *testing.T) {
	rr := httptest.NewRecorder()

	payload := JSONResponse{
		Error:   false,
		Message: "Foo",
	}

	headers := make(http.Header)
	headers.Add("Foo", "Bar")

	err := testApp.writeJSON(rr, http.StatusOK, payload, headers)
	if err != nil {
		t.Errorf("failed to write JSON: %v", err)
	}

	testApp.Environment = "production"
	err = testApp.writeJSON(rr, http.StatusOK, payload, headers)
	if err != nil {
		t.Errorf("failed to write JSON in production env: %v", err)
	}

	testApp.Environment = "development"
}
