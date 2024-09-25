package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Data struct {
	Users    map[string]string `json:"users"`
	Switches map[string]Switch `json:"switches"`
}

var data = Data{
	Users:    make(map[string]string),
	Switches: make(map[string]Switch),
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

	if fn == "" {

	}

	if Config.AppEnv == "development" {
		fmt.Println(data)
	}

	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(fn, file, 0644)
}
