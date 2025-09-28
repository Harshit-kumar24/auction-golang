package utils

import (
	"encoding/json"
	"log"
	"os"
)

func WriteJSONToFile(fileName string, data interface{}) error {

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	// Write JSON bytes to file
	err = os.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		return err
	}
	log.Println("successfully wrote to file...")
	return nil
}
