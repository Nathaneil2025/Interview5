package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func setupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/health", healthHandler).Methods("GET")
	r.HandleFunc("/notifications", getNotificationsHandler).Methods("GET")
	r.HandleFunc("/notifications/{id}", getNotificationHandler).Methods("GET")
	r.HandleFunc("/notifications", createNotificationHandler).Methods("POST")
	return r
}

func TestHealthHandler(t *testing.T) {
	router := setupRouter()

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response HealthResponse
	json.Unmarshal(rr.Body.Bytes(), &response)

	if response.Status != "healthy" {
		t.Errorf("handler returned wrong status: got %v want %v", response.Status, "healthy")
	}

	if response.Service != "notification-service" {
		t.Errorf("handler returned wrong service: got %v want %v", response.Service, "notification-service")
	}
}

func TestGetNotificationsHandler(t *testing.T) {
	router := setupRouter()

	req, err := http.NewRequest("GET", "/notifications", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response []Notification
	json.Unmarshal(rr.Body.Bytes(), &response)

	if len(response) == 0 {
		t.Error("handler returned empty notifications list")
	}
}

func TestCreateNotificationHandler(t *testing.T) {
	router := setupRouter()

	input := NotificationCreate{
		UserID:  1,
		Type:    "email",
		Message: "Test notification",
	}
	body, _ := json.Marshal(input)

	req, err := http.NewRequest("POST", "/notifications", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var response Notification
	json.Unmarshal(rr.Body.Bytes(), &response)

	if response.Message != "Test notification" {
		t.Errorf("handler returned wrong message: got %v want %v", response.Message, "Test notification")
	}
}

func TestCreateNotificationHandlerEmptyMessage(t *testing.T) {
	router := setupRouter()

	input := NotificationCreate{
		UserID:  1,
		Type:    "email",
		Message: "",
	}
	body, _ := json.Marshal(input)

	req, err := http.NewRequest("POST", "/notifications", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}
