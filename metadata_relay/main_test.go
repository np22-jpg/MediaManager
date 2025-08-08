package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"relay/app"
)

func TestRootHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.RootHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"message":"metadata_relay API. See docs for details."}`
	if rr.Body.String() != expected+"\n" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestRouter(t *testing.T) {
	router := app.NewRouter()
	router.GET("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("test"))
	})

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("router returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if rr.Body.String() != "test" {
		t.Errorf("router returned unexpected body: got %v want %v",
			rr.Body.String(), "test")
	}
}
