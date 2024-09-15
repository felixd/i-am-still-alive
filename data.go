package main

import (
	"encoding/json"
	"os"
)

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
