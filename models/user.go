package models

import (
    "database/sql"
    "errors"
    "dead-man-switch/internal/db"
)

type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

func CreateUser(user *User) error {
    query := "INSERT INTO users (username, email, password) VALUES (?, ?, ?)"
    _, err := db.DB.Exec(query, user.Username, user.Email, user.Password)
    return err
}

func GetUserByUsername(username string) (*User, error) {
    var user User
    query := "SELECT id, username, email, password FROM users WHERE username = ?"
    row := db.DB.QueryRow(query, username)
    err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
    if err == sql.ErrNoRows {
        return nil, errors.New("user not found")
    }
    return &user, err
}
