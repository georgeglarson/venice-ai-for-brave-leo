package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// backupPreferencesFile creates a backup of the Preferences file
func backupPreferencesFile(filePath string) (string, error) {
	backupPath := fmt.Sprintf("%s.backup_%s", filePath, time.Now().Format("20060102_150405"))
	input, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	err = os.WriteFile(backupPath, input, 0644)
	if err != nil {
		return "", err
	}
	return backupPath, nil
}

// readPreferencesFile reads and parses the Preferences file
func readPreferencesFile(filePath string) (*BravePreferences, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var preferences BravePreferences
	err = json.Unmarshal(data, &preferences)
	if err != nil {
		return nil, err
	}
	return &preferences, nil
}

// ensureStructureExists ensures that the required JSON structure exists
func ensureStructureExists(preferences *BravePreferences) {
	if preferences.Brave.AIChat.CustomModels == nil {
		preferences.Brave.AIChat.CustomModels = []CustomModel{}
	}
}

// writePreferencesFile writes the updated preferences back to the file
func writePreferencesFile(filePath string, preferences *BravePreferences) error {
	originalData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	var originalJSON map[string]interface{}
	err = json.Unmarshal(originalData, &originalJSON)
	if err != nil {
		return err
	}

	if _, ok := originalJSON["brave"]; !ok {
		originalJSON["brave"] = map[string]interface{}{}
	}
	brave := originalJSON["brave"].(map[string]interface{})
	if _, ok := brave["ai_chat"]; !ok {
		brave["ai_chat"] = map[string]interface{}{}
	}
	aiChat := brave["ai_chat"].(map[string]interface{})
	
	customModelsJSON, err := json.Marshal(preferences.Brave.AIChat.CustomModels)
	if err != nil {
		return err
	}
	var customModels []interface{}
	err = json.Unmarshal(customModelsJSON, &customModels)
	if err != nil {
		return err
	}
	aiChat["custom_models"] = customModels
	aiChat["default_model_key"] = preferences.Brave.AIChat.DefaultModelKey

	updatedData, err := json.MarshalIndent(originalJSON, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, updatedData, 0644)
}

// restoreFromBackup restores the Preferences file from backup
func restoreFromBackup(backupFile, originalFile string) {
	fmt.Println("Restoring from backup...")
	input, err := os.ReadFile(backupFile)
	if err != nil {
		fmt.Printf("Error reading backup file: %v\n", err)
		return
	}
	err = os.WriteFile(originalFile, input, 0644)
	if err != nil {
		fmt.Printf("Error restoring from backup: %v\n", err)
		return
	}
	fmt.Println("Successfully restored from backup.")
}