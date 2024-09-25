package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Healthy!"})
}

func HashPassword(p string) string {
	sha256.Sum256([]byte(p))
	return p
}

func Signup(c *gin.Context) {

	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, exists := data.Users[user.Username]; exists {
		c.JSON(http.StatusConflict, gin.H{"message": "User already exists"})
		return
	}

	data.Users[user.Username] = HashPassword(user.Password)
	if err := SaveData(Config.DataFile); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Signup successful"})
}

func CheckinToken(c *gin.Context) {
	token := c.Param("token")
	log.Printf("Received token: %s", token)

	// Perform token-related logic here

	c.JSON(http.StatusOK, gin.H{
		"message": "Token received",
		"token":   token,
	})
}

// Generate Token for updating switch without Authorization
// r.GET("/checkin/:token", CheckinToken)
func SwitchGenerateToken(c *gin.Context) {

}

func Login(c *gin.Context) {
	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if pwd, exists := data.Users[user.Username]; !exists || pwd != HashPassword(user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "No authorization token provided"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method")
			}
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		if err != nil {
			fmt.Errorf("Error parsing token")
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("username", claims["username"])
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			c.Abort()
		}
	}
}
