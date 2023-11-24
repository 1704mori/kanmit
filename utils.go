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

// func WriteJSONToFile(filename string, data interface{}) error {
// 	jsonData, err := json.Marshal(data)
// 	if err != nil {
// 		return err
// 	}

// 	file, err := os.Create(filename)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	_, err = file.Write(jsonData)
// 	return err
// }

func WriteJSONToFile(filename string, key string, value interface{}) error {
	// Read existing JSON data from the file
	existingData, err := os.ReadFile(filename)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Unmarshal existing data into a map
	var existingMap map[string]interface{}
	if len(existingData) > 0 {
		err = json.Unmarshal(existingData, &existingMap)
		if err != nil {
			return err
		}
	} else {
		// If the file doesn't exist or is empty, create a new map
		existingMap = make(map[string]interface{})
	}

	// Update or add the key-value pair
	existingMap[key] = value

	// Marshal the updated data
	updatedData, err := json.Marshal(existingMap)
	if err != nil {
		return err
	}

	// Write the updated data back to the file
	err = os.WriteFile(filename, updatedData, 0644)
	if err != nil {
		return err
	}

	return nil
}
