package gohttp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	c "../controller"
)

func TestHealthCheckHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v1/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(c.Health)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got '%v' want '%v'",
			status, http.StatusOK)
	}

	expected := `{"data":"OK"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got '%v' want '%v'",
			rr.Body.String(), expected)
	}
}
