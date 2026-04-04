package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func TestGETHelloName_OK(t *testing.T) {
	r := newRouter()
	req := httptest.NewRequest(http.MethodGet, "/hello/Alice", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d, body: %s", http.StatusOK, rec.Code, rec.Body.String())
	}
	var body struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("json: %v", err)
	}
	if body.Message != "hello there Alice" {
		t.Fatalf("message: want %q, got %q", "hello there Alice", body.Message)
	}
}

func TestGETHelloName_ValidationError(t *testing.T) {
	r := newRouter()
	req := httptest.NewRequest(http.MethodGet, "/hello/x", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want %d, got %d, body: %s", http.StatusBadRequest, rec.Code, rec.Body.String())
	}
	var body struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("json: %v", err)
	}
	if body.Error == "" {
		t.Fatal("expected non-empty error field")
	}
}

func TestPOSTProfile_Created(t *testing.T) {
	r := newRouter()
	payload := map[string]any{
		"display_name": "Timur",
		"email":        "t@example.com",
		"age":          20,
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/profile", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status: want %d, got %d, body: %s", http.StatusCreated, rec.Code, rec.Body.String())
	}
	var body struct {
		DisplayName string `json:"display_name"`
		Email       string `json:"email"`
		Age         int    `json:"age"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("json: %v", err)
	}
	if body.DisplayName != "Timur" || body.Email != "t@example.com" || body.Age != 20 {
		t.Fatalf("response: %+v", body)
	}
}

func TestPOSTProfile_ValidationError(t *testing.T) {
	r := newRouter()
	payload := map[string]any{
		"display_name": "Timur",
		"email":        "not-an-email",
		"age":          20,
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/profile", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want %d, got %d, body: %s", http.StatusBadRequest, rec.Code, rec.Body.String())
	}
}
