package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

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

	if err := SaveData(Config.DataFile); err != nil {
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

		if err := SaveData(Config.DataFile); err != nil {
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

		if err := SaveData(Config.DataFile); err != nil {
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
				subject := "Your Dead Man Switch has been triggered"
				body := "This is to inform you that your Dead Man Switch has been triggered."

				err := SendEmail(switchData.Recipients, subject, body)
				if err != nil {
					fmt.Printf("Error sending email to %s: %v\n", switchData.User, err)
				} else {
					fmt.Printf("Emails sent to: %v\n", switchData.Recipients)
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
