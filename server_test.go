package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleAbout(t *testing.T) {
	// prepare
	srv := server{
		router: http.NewServeMux(),
		//db:    mockDatabase,
		//email: mockEmailSender,
	}
	srv.routes()
	r := httptest.NewRequest("GET", "/about", nil)
	w := httptest.NewRecorder()
	srv.router.ServeHTTP(w, r)
	res := w.Result()
	// test
	if res.StatusCode != http.StatusOK {
		t.Error("wrong status code")
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("%s at response body", err)
	}
	s := string(body)
	if s != "I am about!" {
		t.Error("response is not correct")
	}
}

func TestHandleIndex(t *testing.T) {
	// prepare
	srv := server{
		router: http.NewServeMux(),
		//db:    mockDatabase,
		//email: mockEmailSender,
	}
	srv.routes()
	// test root dir
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	srv.router.ServeHTTP(w, r)
	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Error("wrong status code")
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("%s response body", err)
	}
	s := string(body)
	if s != "This is index! Passed data is tsumutsumu." {
		t.Error("response is not correct")
	}
	// test not matched dir
	r = httptest.NewRequest("GET", "/yeah", nil)
	w = httptest.NewRecorder()
	srv.router.ServeHTTP(w, r)
	res = w.Result()
	if res.StatusCode != http.StatusOK {
		t.Error("wrong status code")
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("%s at response body", err)
	}
	s = string(body)
	if s != "This is index! Passed data is tsumutsumu." {
		t.Error("response is not correct")
	}
}

func TestHandleAPI(t *testing.T) {
	// prepare
	srv := server{
		router: http.NewServeMux(),
		//db:    mockDatabase,
		//email: mockEmailSender,
	}
	srv.routes()
	type request struct {
		Name string
	}
	type response struct {
		Greeting string `json:"greetinig"`
	}
	// test normal name
	req := request{Name: "john"}
	resp := response{}
	r := httptest.NewRequest("GET", "/api/?name="+req.Name, nil)
	w := httptest.NewRecorder()
	srv.router.ServeHTTP(w, r)
	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Error("wrong status code")
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("%s at response body", err)
	}
	if err = json.Unmarshal(body, &resp); err != nil {
		t.Errorf("%s at response json parse", err)
	}
	if resp.Greeting != "Hello, "+req.Name {
		t.Error("response json is not correct")
	}
	// test specific name
	req = request{Name: "P"}
	resp = response{}
	r = httptest.NewRequest("GET", "/api/?name="+req.Name, nil)
	w = httptest.NewRecorder()
	srv.router.ServeHTTP(w, r)
	res = w.Result()
	if res.StatusCode != http.StatusOK {
		t.Error("wrong status code")
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("%s at response body", err)
	}
	if err = json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("%s at response json parse", err)
	}
	if resp.Greeting != "bakananodesuka, "+req.Name {
		t.Error("response json is not correct")
	}
	// test no name
	resp = response{}
	r = httptest.NewRequest("GET", "/api/", nil)
	w = httptest.NewRecorder()
	srv.router.ServeHTTP(w, r)
	res = w.Result()
	if res.StatusCode != http.StatusOK {
		t.Error("wrong status code")
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("%s at response body", err)
	}
	if err = json.Unmarshal(body, &resp); err != nil {
		t.Errorf("%s at response json parse", err)
	}
	if resp.Greeting != "Hello, This is API." {
		t.Error("response json is not correct")
	}
}

// other tests are abbreviated
