package main

import (
	"encoding/json"
	"os"
	"time"
)

type DeadManSwitch struct {
	User        string        `json:"user"`
	SwitchDelay time.Duration `json:"switch_delay"`
	TriggerAt   time.Time     `json:"trigger_at"`
	Recipients  []string      `json:"recipients"`
	Message     []string      `json:"message"`
}

type Data struct {
	Users    map[string]string        `json:"users"`
	Switches map[string]DeadManSwitch `json:"switches"`
}

var data = Data{
	Users:    make(map[string]string),
	Switches: make(map[string]DeadManSwitch),
}

func LoadData(fn string) error {
	file, err := os.ReadFile(fn)
	if err != nil {
		if os.IsNotExist(err) {
			return SaveData(fn)
		}
		return err
	}
	return json.Unmarshal(file, &data)
}

func SaveData(fn string) error {
	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(fn, file, 0644)
}
