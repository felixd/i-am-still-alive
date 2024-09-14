package handlers

import (
    "encoding/json"
    "dead-man-switch/models"
    "net/http"
    "log"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
        return
    }

    var user models.User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    err = models.CreateUser(&user)
    if err != nil {
        http.Error(w, "Failed to create user", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
        return
    }

    var user models.User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    existingUser, err := models.GetUserByUsername(user.Username)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    if existingUser.Password != user.Password {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    w.WriteHeader(http.StatusOK)
}
