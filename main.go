package main

/*
 * I am still alive - Dead Person Switch
 * (c) FlameIT - Immersion Cooling - Paweł Wojciechowski
 */

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type User struct {
	ID        string        `json:"id"`
	Password  string        `json:"password"`
	LastCheck time.Time     `json:"last_check"`
	Message   string        `json:"message"`
	Timeout   time.Duration `json:"timeout"`
	IsAlive   bool          `json:"is_alive"`
}

type DeadManSwitch struct {
	mu            sync.Mutex
	users         map[string]*User
	checkInterval time.Duration
}

func NewDeadManSwitch(checkInterval time.Duration) *DeadManSwitch {
	return &DeadManSwitch{
		users:         make(map[string]*User),
		checkInterval: checkInterval,
	}
}

// Wczytaj użytkowników z pliku JSON
func (dms *DeadManSwitch) LoadUsers(filename string) error {
	dms.mu.Lock()
	defer dms.mu.Unlock()

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &dms.users)
	if err != nil {
		return err
	}
	return nil
}

// Zapisz użytkowników do pliku JSON
func (dms *DeadManSwitch) SaveUsers(filename string) error {
	dms.mu.Lock()
	defer dms.mu.Unlock()

	file, err := json.MarshalIndent(dms.users, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, file, 0644)
}

// Dodanie nowego użytkownika
func (dms *DeadManSwitch) AddUser(id string, password string, timeout time.Duration, message string) {
	dms.mu.Lock()
	defer dms.mu.Unlock()
	dms.users[id] = &User{
		ID:        id,
		Password:  password,
		LastCheck: time.Now(),
		Message:   message,
		Timeout:   timeout,
		IsAlive:   true,
	}
}

// Zmiana wiadomości użytkownika
func (dms *DeadManSwitch) UpdateMessage(id string, message string) {
	dms.mu.Lock()
	defer dms.mu.Unlock()
	if user, exists := dms.users[id]; exists {
		user.Message = message
	}
}

// Odbieranie sygnału KeepAlive od użytkownika
func (dms *DeadManSwitch) KeepAlive(id string) {
	dms.mu.Lock()
	defer dms.mu.Unlock()
	if user, exists := dms.users[id]; exists {
		user.LastCheck = time.Now()
		user.IsAlive = true
		fmt.Printf("Keepalive received from user: %s\n", id)
	} else {
		fmt.Printf("User %s not found\n", id)
	}
}

// Sprawdzanie stanu użytkowników
func (dms *DeadManSwitch) CheckUsers() {
	for {
		time.Sleep(dms.checkInterval)
		dms.mu.Lock()
		for _, user := range dms.users {
			if time.Since(user.LastCheck) > user.Timeout {
				user.IsAlive = false
				fmt.Printf("Dead Man's Switch triggered for user %s! Sending message: %s\n", user.ID, user.Message)
				// Tutaj można dodać logikę wysyłania powiadomienia, np. e-mail.
			} else {
				fmt.Printf("User %s is alive.\n", user.ID)
			}
		}
		dms.mu.Unlock()
	}
}

// Autoryzacja użytkownika (Basic Auth)
func basicAuth(next http.HandlerFunc, dms *DeadManSwitch) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || !dms.ValidateUser(user, pass) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Please enter username and password"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// Weryfikacja użytkownika
func (dms *DeadManSwitch) ValidateUser(id, password string) bool {
	dms.mu.Lock()
	defer dms.mu.Unlock()
	user, exists := dms.users[id]
	return exists && user.Password == password
}

func main() {
	// Konfiguracja Dead Man's Switch
	checkInterval := 5 * time.Second // Sprawdzaj co 5 sekund
	dms := NewDeadManSwitch(checkInterval)

	// Wczytaj użytkowników z pliku
	userFile := "users.json"
	if _, err := os.Stat(userFile); err == nil {
		err := dms.LoadUsers(userFile)
		if err != nil {
			log.Fatalf("Error loading users: %v", err)
		}
	} else {
		fmt.Println("No user file found, starting with empty user list.")
	}

	// Start sprawdzania użytkowników
	go dms.CheckUsers()

	// Endpoint HTTP do potwierdzania "życia" dla użytkowników
	http.HandleFunc("/keepalive", basicAuth(func(w http.ResponseWriter, r *http.Request) {
		userID, _, _ := r.BasicAuth()
		dms.KeepAlive(userID)
		fmt.Fprintf(w, "Keepalive received for user: %s", userID)
	}, dms))

	// Endpoint do dodawania użytkowników
	http.HandleFunc("/adduser", basicAuth(func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user")
		password := r.URL.Query().Get("password")
		message := r.URL.Query().Get("message")
		timeoutStr := r.URL.Query().Get("timeout")

		if userID == "" || password == "" || message == "" || timeoutStr == "" {
			http.Error(w, "Missing parameters", http.StatusBadRequest)
			return
		}

		timeout, err := time.ParseDuration(timeoutStr)
		if err != nil {
			http.Error(w, "Invalid timeout value", http.StatusBadRequest)
			return
		}

		dms.AddUser(userID, password, timeout, message)
		err = dms.SaveUsers(userFile)
		if err != nil {
			http.Error(w, "Error saving user data", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "User %s added with message: %s and timeout: %s", userID, message, timeoutStr)
	}, dms))

	// Endpoint do zmiany wiadomości użytkownika
	http.HandleFunc("/updatemessage", basicAuth(func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user")
		message := r.URL.Query().Get("message")

		if userID == "" || message == "" {
			http.Error(w, "User ID and message are required", http.StatusBadRequest)
			return
		}

		dms.UpdateMessage(userID, message)
		err := dms.SaveUsers(userFile)
		if err != nil {
			http.Error(w, "Error saving user data", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Message updated for user: %s", userID)
	}, dms))

	// Uruchomienie serwera
	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
