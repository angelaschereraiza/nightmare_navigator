package telegram_bot

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

const (
	downloadDir  = "data"
	jsonFilename = "user_ids.json"
)

var userIdsFile = filepath.Join(downloadDir, jsonFilename)

// Function to load the userId list from a JSON file
func loadUserIds() []int64 {
	var userIds []int64

	// Tries to read json file
	jsonData, err := os.ReadFile(userIdsFile)
	if err != nil {
		log.Println("Error loading userIds:", err)
	}

	err = json.Unmarshal(jsonData, &userIds)
	if err != nil {
		log.Println("Error loading userIds:", err)
	}

	return userIds
}

func saveUserIds(userId int64) {
	userIds := loadUserIds()
	if userIds == nil {
		// Create the directory if it does not exist
		err := os.MkdirAll(downloadDir, 0755)
		if err != nil {
			log.Println("Error creating directory:", err)
		}

		// Creates a new json file if it does not exist
		jsonData, err := os.Create(downloadDir + "/" + jsonFilename)
		if err != nil {
			log.Println("Error creating json file:", err)
		}
		defer jsonData.Close()
	}

	userIds = append(userIds, userId)

	jsonData, err := json.MarshalIndent(userIds, "", "    ")
	if err != nil {
		log.Println("Error saving userIds file:", err)
	}

	err = os.WriteFile(userIdsFile, jsonData, 0644)
	if err != nil {
		log.Println("Error saving userIds file:", err)
	}
}
