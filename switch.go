package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Switch struct {
	Id         int           `json:"id"`
	Token      string        `json:"token"`
	User       string        `json:"user"`
	Duration   time.Duration `json:"duration"`
	TriggerAt  time.Time     `json:"trigger_at"`
	Recipients []string      `json:"recipients"`
	Message    string        `json:"message"`
}

func SwitchCreate(c *gin.Context) {

	r := Switch{}
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username, _ := c.Get("username")
	triggerTime := time.Now().Add(time.Duration(r.Duration) * time.Hour)
	r.User = username.(string)
	r.TriggerAt = triggerTime

	data.Switches[username.(string)] = r

	// Save to "db" (JSON file)
	if err := SaveData(Config.DataFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Switch created", "trigger_at": triggerTime})
}

func SwitchUpdate(c *gin.Context) {
	username, _ := c.Get("username")

	r := Switch{}
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if switchData, exists := data.Switches[username.(string)]; exists {
		if r.Duration > 0 {
			d := time.Duration(r.Duration) * time.Hour
			switchData.TriggerAt = time.Now().Add(d)
			switchData.Duration = d
		}
		if len(r.Recipients) > 0 {
			switchData.Recipients = r.Recipients
		}
		data.Switches[username.(string)] = switchData

		if err := SaveData(Config.DataFile); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Switch updated", "switch": switchData})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"message": "No switch found"})
	}
}

func SwitchCheckin(c *gin.Context) {
	username, _ := c.Get("username")
	if switchData, exists := data.Switches[username.(string)]; exists {
		switchData.TriggerAt = time.Now().Add(time.Hour * switchData.Duration)
		data.Switches[username.(string)] = switchData

		if err := SaveData(Config.DataFile); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Check-in successful", "next_trigger": switchData.TriggerAt})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"message": "No switch found"})
	}
}

func SwitchDelete(c *gin.Context) {
	username, _ := c.Get("username")
	delete(data.Switches, username.(string))

	if err := SaveData(Config.DataFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Switch deleted"})
}

func CheckExpiredSwitches() {
	for {
		for username, switchData := range data.Switches {
			if time.Now().After(switchData.TriggerAt) {
				subject := "Dead Man Switch of '" + username + "' has been triggered"
				body := "This is to inform you that Dead Man Switch of '" + username + "' has been triggered.\n"
				body += "Message is contained between ### lines: \n"
				body += "\n######################\n"
				body += switchData.Message
				body += "\n######################\n"

				err := SendEmail(switchData.Recipients, subject, body)
				if err != nil {
					log.Printf("Error sending email to %s: %v\n", switchData.User, err)
				} else {
					log.Printf("Email sent to: %v\n", switchData.Recipients)
				}

				delete(data.Switches, username)

				if err := SaveData(Config.DataFile); err != nil {
					fmt.Println("Error saving data after trigger")
				}
			}
		}
		time.Sleep(time.Minute * 1)
	}
}
