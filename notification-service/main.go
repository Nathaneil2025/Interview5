package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

type Notification struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	Read      bool      `json:"read"`
}

type NotificationCreate struct {
	UserID  int    `json:"user_id"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

var (
	notifications = []Notification{
		{ID: 1, UserID: 1, Type: "email", Message: "Welcome to our service!", CreatedAt: time.Now(), Read: false},
		{ID: 2, UserID: 1, Type: "sms", Message: "Your transaction was successful", CreatedAt: time.Now(), Read: true},
	}
	mu     sync.Mutex
	nextID = 3
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:  "healthy",
		Service: "notification-service",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

func getNotificationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	for _, n := range notifications {
		if string(rune(n.ID+'0')) == id || (n.ID >= 10 && string(n.ID) == id) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(n)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Notification not found"})
}

func createNotificationHandler(w http.ResponseWriter, r *http.Request) {
	var input NotificationCreate
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	if input.Message == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Message is required"})
		return
	}

	mu.Lock()
	notification := Notification{
		ID:        nextID,
		UserID:    input.UserID,
		Type:      input.Type,
		Message:   input.Message,
		CreatedAt: time.Now(),
		Read:      false,
	}
	nextID++
	notifications = append(notifications, notification)
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(notification)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/health", healthHandler).Methods("GET")
	r.HandleFunc("/notifications", getNotificationsHandler).Methods("GET")
	r.HandleFunc("/notifications/{id}", getNotificationHandler).Methods("GET")
	r.HandleFunc("/notifications", createNotificationHandler).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Notification service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
