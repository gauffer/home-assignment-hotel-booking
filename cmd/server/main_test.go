package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func testApp() *chi.Mux {
	r := SetupMain(nil)
	r = SetupOrders(r)
	return r
}

func TestPing(t *testing.T) {
	req := httptest.NewRequest("GET", "/ping", nil)
	respRecorder := httptest.NewRecorder()

	testApp().ServeHTTP(respRecorder, req)

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

	testApp().ServeHTTP(respRecorder, req)

	if respRecorder.Code != http.StatusCreated {
		t.Errorf(
			"Expected status code %d, got %d",
			http.StatusCreated,
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

	testApp().ServeHTTP(respRecorder, req)

	if respRecorder.Code != http.StatusBadRequest {
		t.Errorf(
			"Expected status code %d, got %d",
			http.StatusCreated,
			respRecorder.Code,
		)
	}

	responseString := strings.TrimSpace(respRecorder.Body.String())
	if !strings.Contains(responseString, "hotel_id is required") {
		t.Error("Expected 'hotel_id is required' in the response body")
	}
}
