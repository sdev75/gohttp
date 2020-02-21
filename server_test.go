package gohttp

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	cfg = Config{
		Addr:              "0.0.0.0:9999",
		ReadTimeout:       time.Second * 5,
		ReadHeaderTimeout: time.Second * 5,
		WriteTimeout:      time.Second * 30,
		IdleTimeout:       time.Second * 30,
	}
	router = NewSimpleRouter()
)

func TestServer(t *testing.T) {

	router.Match = func(w http.ResponseWriter, r *http.Request) *SimpleRouterMatch {
		res := &SimpleRouterMatch{}
		res.Handler = func(w http.ResponseWriter, r *http.Request) {
			WriteSuccess(r.Context(), w, "OK")
		}
		return res
	}

	Init(router, cfg, nil)
	go Start()
	defer Stop()

	req, err := http.NewRequest("GET", "health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	m := router.Match(rr, req)
	m.Handler(rr, req)

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
