package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func mockedApp() *chi.Mux {
	r := setupChi(nil)
	r = setupOrdersRoute(r)
	return r
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func TestPing(t *testing.T) {
	req := httptest.NewRequest("GET", "/ping", nil)
	respRecorder := httptest.NewRecorder()

	mockedApp().ServeHTTP(respRecorder, req)

	if respRecorder.Code != http.StatusOK {
		t.Errorf(
			"Expected status code %d, got %d",
			http.StatusOK,
			respRecorder.Code,
		)
	}

	body := strings.TrimSpace(respRecorder.Body.String())
	if body != "." {
		t.Errorf("Expected response '.', got '%s'", body)
	}
}

func TestCreateOrder(t *testing.T) {
	payload := []byte(`{
	"hotel_id": "reddison",
	"room_id": "lux",
	"email": "guest@mail.ru",
	"from": "2024-01-02T00:00:00Z",
	"to": "2024-01-04T00:00:00Z"
	}`)

	req := httptest.NewRequest("POST", "/orders", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	respRecorder := httptest.NewRecorder()

	mockedApp().ServeHTTP(respRecorder, req)

	if respRecorder.Code != http.StatusOK {
		t.Errorf(
			"Expected status code %d, got %d",
			http.StatusOK,
			respRecorder.Code,
		)
	}
}

func TestCreateInvalidOrder(t *testing.T) {
	payload := []byte(`{
	"room_id": "lux",
	"email": "guest@mail.ru",
	"from": "2024-01-02T00:00:00Z",
	"to": "2024-01-04T00:00:00Z"
	}`)

	req := httptest.NewRequest("POST", "/orders", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	respRecorder := httptest.NewRecorder()

	mockedApp().ServeHTTP(respRecorder, req)

	if respRecorder.Code != http.StatusUnprocessableEntity {
		t.Errorf(
			"Expected status code %d, got %d",
			http.StatusUnprocessableEntity, // Updated status expectation
			respRecorder.Code,
		)
	}

	// Decode error response
	var errResponse ErrorResponse
	err := json.Unmarshal(respRecorder.Body.Bytes(), &errResponse)
	if err != nil {
		t.Fatalf("Failed to decode error response: %v", err)
	}

	expectedMessage := "Missing required fields: hotel_id"
	if errResponse.Message != expectedMessage {
		t.Errorf(
			"Expected %s in error message, got: %s",
			expectedMessage,
			errResponse.Message,
		)
	}
}
