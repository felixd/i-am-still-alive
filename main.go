package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const ()

type Config struct {
	dataFile     string "data.json"
	smtpHost     string "smtp.gmail.com"
	smtpPort     string "587"
	senderEmail  string "your-email@gmail.com"
	senderPass   string "your-email-password"
	jwtSecretKey string "secret_key"
}

type DeadManSwitch struct {
	User       string    `json:"user"`
	TriggerAt  time.Time `json:"trigger_at"`
	Recipients []string  `json:"recipients"`
	Message    []string  `json:"message"`
}

type Data struct {
	Users    map[string]string        `json:"users"`
	Switches map[string]DeadManSwitch `json:"switches"`
}

var data = Data{
	Users:    make(map[string]string),
	Switches: make(map[string]DeadManSwitch),
}

func main() {
	r := gin.Default()

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load data from JSON file
	if err := LoadData(); err != nil {
		panic(err)
	}

	// Routes
	r.POST("/signup", Signup)
	r.POST("/login", Login)

	authorized := r.Group("/switch")
	authorized.Use(AuthMiddleware())
	{
		authorized.POST("/create", CreateSwitch)
		authorized.GET("/checkin", CheckinSwitch)
		authorized.DELETE("/delete", DeleteSwitch)
		authorized.PUT("/update", UpdateSwitch) // Endpoint for updating the switch
	}

	// Periodic check for triggered switches
	go CheckExpiredSwitches()

	r.Run(":8080")
}

func LoadData() error {
	file, err := os.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return SaveData()
		}
		return err
	}
	return json.Unmarshal(file, &data)
}

func SaveData() error {
	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, file, 0644)
}

func SendEmail(recipients []string, subject, body string) error {
	auth := smtp.PlainAuth("", senderEmail, senderPass, smtpHost)
	to := strings.Join(recipients, ", ")
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body + "\r\n")
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, recipients, msg)
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

	data.Users[user.Username] = user.Password
	if err := SaveData(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Signup successful"})
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

	if pwd, exists := data.Users[user.Username]; !exists || pwd != user.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecretKey))
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
			return []byte(jwtSecretKey), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
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

func CreateSwitch(c *gin.Context) {
	username, _ := c.Get("username")
	var switchRequest struct {
		Duration   int      `json:"duration"`
		Recipients []string `json:"recipients"`
	}

	if err := c.ShouldBindJSON(&switchRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	triggerTime := time.Now().Add(time.Duration(switchRequest.Duration) * time.Hour)
	data.Switches[username.(string)] = DeadManSwitch{
		User:       username.(string),
		TriggerAt:  triggerTime,
		Recipients: switchRequest.Recipients,
	}

	if err := SaveData(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Switch created", "trigger_at": triggerTime})
}

func UpdateSwitch(c *gin.Context) {
	username, _ := c.Get("username")
	var updateRequest struct {
		Duration   int      `json:"duration"`
		Recipients []string `json:"recipients"`
	}

	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if switchData, exists := data.Switches[username.(string)]; exists {
		if updateRequest.Duration > 0 {
			switchData.TriggerAt = time.Now().Add(time.Duration(updateRequest.Duration) * time.Hour)
		}
		if len(updateRequest.Recipients) > 0 {
			switchData.Recipients = updateRequest.Recipients
		}
		data.Switches[username.(string)] = switchData

		if err := SaveData(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Switch updated", "switch": switchData})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"message": "No switch found"})
	}
}

func CheckinSwitch(c *gin.Context) {
	username, _ := c.Get("username")
	if switchData, exists := data.Switches[username.(string)]; exists {
		switchData.TriggerAt = time.Now().Add(time.Hour * 24)
		data.Switches[username.(string)] = switchData

		if err := SaveData(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Check-in successful", "next_trigger": switchData.TriggerAt})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"message": "No switch found"})
	}
}

func DeleteSwitch(c *gin.Context) {
	username, _ := c.Get("username")
	delete(data.Switches, username.(string))

	if err := SaveData(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Switch deleted"})
}

func CheckExpiredSwitches() {
	for {
		for username, switchData := range data.Switches {
			if time.Now().After(switchData.TriggerAt) {
				subject := "Your Dead Man Switch has been triggered"
				body := "This is to inform you that your Dead Man Switch has been triggered."

				err := SendEmail(switchData.Recipients, subject, body)
				if err != nil {
					fmt.Printf("Error sending email to %s: %v\n", switchData.User, err)
				} else {
					fmt.Printf("Emails sent to: %v\n", switchData.Recipients)
				}

				delete(data.Switches, username)

				if err := SaveData(); err != nil {
					fmt.Println("Error saving data after trigger")
				}
			}
		}
		time.Sleep(time.Minute * 1)
	}
}
