package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestHelloHandler(t *testing.T) {
    // Buat request GET ke endpoint
    req := httptest.NewRequest(http.MethodGet, "/api/hello", nil)
    // Rekam response dari handler
    rr := httptest.NewRecorder()

    HelloHandler(rr, req)
    res := rr.Result()
    defer res.Body.Close()

    // Cek status code
    if res.StatusCode != http.StatusOK {
        t.Errorf("Status code salah: dapat %v, mau %v", res.StatusCode, http.StatusOK)
    }

    // Cek Content-Type header
    if ct := res.Header.Get("Content-Type"); ct != "application/json" {
        t.Errorf("Content-Type salah: dapat %q, mau %q", ct, "application/json")
    }

    // Cek body JSON
    expected := `{"Text":"Halo dari Golang!"}` + "\n"
    body := rr.Body.String()
    if body != expected {
        t.Errorf("Body salah: dapat %q, mau %q", body, expected)
    }
}