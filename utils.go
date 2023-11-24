package main

import (
	"encoding/json"
	"os"
)

func ReadJSONFromFile(filename string, v interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(v)
}

func WriteJSONToFile(filename string, key string, value interface{}) error {
	existingData, err := os.ReadFile(filename)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var existingMap map[string]interface{}
	if len(existingData) > 0 {
		err = json.Unmarshal(existingData, &existingMap)
		if err != nil {
			return err
		}
	} else {
		existingMap = make(map[string]interface{})
	}

	existingMap[key] = value

	updatedData, err := json.Marshal(existingMap)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, updatedData, 0644)
	if err != nil {
		return err
	}

	return nil
}
