package tests

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"testing"

	"github.com/mikeyQwn/greq"
)

const (
	TESTING_ADDR = "127.0.0.1:7128"
	REQ_ADDR     = "http://" + TESTING_ADDR

	JSON_ENDPOINT   = "/json_endpoint"
	STRING_ENDPOINT = "/string_endpoint"
	AUTHED_ENDPOINT = "/authed_endpoint"
)

func mapHandlers(mux *http.ServeMux) {
	// JSON_ENDPOINT
	mux.HandleFunc("POST "+JSON_ENDPOINT, func(w http.ResponseWriter, r *http.Request) {
		data := map[string]any{
			"status":  "ok",
			"value":   10,
			"boolean": true,
			"foo":     []int{1, 2, 3},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	})

	// STRING_ENDPOINT
	mux.HandleFunc("GET "+STRING_ENDPOINT, func(w http.ResponseWriter, r *http.Request) {
		data := []byte("Hello, friend")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})

	//AUTHED_ENDPOINT
	mux.HandleFunc("GET "+AUTHED_ENDPOINT, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Cookie") != "Auth=abc" {
			w.WriteHeader(http.StatusUnauthorized)
			err := map[string]string{"error": "invalid cookie"}
			json.NewEncoder(w).Encode(err)
			return
		}
		if r.Header.Get("Fingerprint") != "fingerprint" {
			w.WriteHeader(http.StatusUnauthorized)
			err := map[string]string{"error": "invalid fingerprint"}
			json.NewEncoder(w).Encode(err)
			return
		}
		w.WriteHeader(http.StatusOK)
		resp := map[string]string{"status": "logged in"}
		json.NewEncoder(w).Encode(resp)
	})
}

func createTestingServer(done chan<- struct{}) {
	mux := http.NewServeMux()
	mapHandlers(mux)
	ln, err := net.Listen("tcp", TESTING_ADDR)
	if err != nil {
		log.Fatal(err)
	}
	done <- struct{}{}
	http.Serve(ln, mux)
}

func init() {
	done := make(chan struct{})
	go createTestingServer(done)
	<-done
}

func TestJsonPost(t *testing.T) {
	type Response struct {
		Status  string
		Value   int
		Boolean bool
		Foo     []int
	}
	resp, err := greq.NewRequest[Response]().MustPost(REQ_ADDR + JSON_ENDPOINT).BaseType()
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status != "ok" {
		t.Fatal("Status is not ok")
	}
	if resp.Value != 10 {
		t.Fatal("Value is not 10")
	}
	if !resp.Boolean {
		t.Fatal("Boolean is not true")
	}
	expectedFoo := []int{1, 2, 3}
	if len(resp.Foo) != len(expectedFoo) {
		t.Fatal("Foo len is right")
	}
	for i, v := range resp.Foo {
		if expectedFoo[i] != v {
			t.Fatal(fmt.Sprintf("Foo [%d] is not %d", i, v))
		}
	}
}

func TestStringGet(t *testing.T) {
	resp := greq.NewRequest[struct{}]().MustGet(REQ_ADDR + STRING_ENDPOINT).String()
	if resp != "Hello, friend" {
		t.Fatal(resp)
	}
}

type ErrorResponse struct {
	Error string
}
type OkResponse struct {
	Status string
}

func TestHeadersA(t *testing.T) {
	resp, err := greq.NewRequest[ErrorResponse]().MustGet(REQ_ADDR + AUTHED_ENDPOINT).BaseType()
	if err != nil {
		t.Fatal(err)
	}
	if resp.Error != "invalid cookie" {
		t.Fatal("Unexpected error:", resp.Error)
	}
}

func TestHeadersB(t *testing.T) {
	resp, err := greq.NewRequest[ErrorResponse]().WithHeader("Cookie", "Auth=abc").MustGet(REQ_ADDR + AUTHED_ENDPOINT).BaseType()
	if err != nil {
		t.Fatal(err)
	}
	if resp.Error != "invalid fingerprint" {
		t.Fatal("Unexpected error:", resp.Error)
	}
}

func TestHeadersC_KV(t *testing.T) {
	resp, err := greq.NewRequest[OkResponse]().WithHeader("Cookie", "Auth=abc").WithHeader("Fingerprint", "fingerprint").MustGet(REQ_ADDR + AUTHED_ENDPOINT).BaseType()
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status != "logged in" {
		t.Fatal("Unexpected status:", resp.Status)
	}
}

func TestHeadersD_map(t *testing.T) {
	resp, err := greq.NewRequest[OkResponse]().WithHeaders(map[string]string{"Cookie": "Auth=abc", "Fingerprint": "fingerprint"}).MustGet(REQ_ADDR + AUTHED_ENDPOINT).BaseType()
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status != "logged in" {
		t.Fatal("Unexpected status:", resp.Status)
	}
}
