package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/google/uuid"
)

func main() {
	// Parse command line arguments
	apiKey := flag.String("key", "", "Venice.AI API key")
	flag.Parse()

	// Check if API key is provided
	if *apiKey == "" {
		fmt.Println("No API key provided via command line, prompting for input...")
		inputKey, ok := GetAPIKeyFromDialog()
		if !ok || inputKey == "" {
			fmt.Println("Error: No API key provided. Operation cancelled.")
			if runtime.GOOS == "windows" {
				ShowMessageBox("Canceled", "No API key provided. Configuration aborted.", MB_OK|MB_ICONINFORMATION)
			}
			os.Exit(1)
		}
		*apiKey = inputKey
		fmt.Println("API key received.")
	}

	fmt.Println("Leo Venice.AI Configuration Tool " + AppVersion)
	fmt.Println("==========================================")

	// Find Brave's Preferences.json file
	preferencesFile, err := findPreferencesFile()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		if runtime.GOOS == "windows" {
			ShowMessageBox("Error", fmt.Sprintf("Could not find Preferences file: %v", err), MB_OK|MB_ICONINFORMATION)
		}
		os.Exit(1)
	}
	fmt.Printf("Found Preferences file at: %s\n", preferencesFile)

	// Check if Brave is running
	if isBraveRunning() {
		fmt.Println("WARNING: Brave browser is currently running.")
		if runtime.GOOS == "windows" {
			ret := ShowMessageBox("Warning", "Brave is running in the background.\n\nFor changes to take effect, Brave must be completely closed.\n\nWould you like to:\n- Click 'Yes' to close Brave and continue\n- Click 'No' to continue without closing Brave\n- Click 'Cancel' to abort", MB_YESNOCANCEL|MB_ICONQUESTION)
			if ret == IDCANCEL {
				fmt.Println("Operation cancelled by user.")
				os.Exit(0)
			} else if ret == IDYES {
				fmt.Println("Attempting to close Brave browser...")
				if !killBraveProcesses() {
					if ShowMessageBox("Warning", "Could not completely close Brave. Changes may not take effect until you restart Brave.\n\nContinue anyway?", MB_OKCANCEL|MB_ICONQUESTION) != IDOK {
						fmt.Println("Operation cancelled by user.")
						os.Exit(0)
					}
				}
			}
		} else {
			fmt.Println("For changes to take effect, Brave must be completely closed.")
			fmt.Print("Would you like to close Brave and continue? (y/n/c) [y=yes, n=continue without closing, c=cancel]: ")
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.ToLower(strings.TrimSpace(response))
			
			if response == "c" {
				fmt.Println("Operation cancelled by user.")
				os.Exit(0)
			} else if response == "y" {
				fmt.Println("Attempting to close Brave browser...")
				if !killBraveProcesses() {
					fmt.Println("WARNING: Could not completely close Brave. Changes may not take effect until you restart Brave.")
					fmt.Print("Continue anyway? (y/n): ")
					reader := bufio.NewReader(os.Stdin)
					response, _ := reader.ReadString('\n')
					if strings.ToLower(strings.TrimSpace(response)) != "y" {
						fmt.Println("Operation cancelled by user.")
						os.Exit(0)
					}
				}
			}
		}
	}

	// Create backup
	backupFile, err := backupPreferencesFile(preferencesFile)
	if err != nil {
		fmt.Printf("Error creating backup: %v\n", err)
		if runtime.GOOS == "windows" {
			ShowMessageBox("Error", fmt.Sprintf("Failed to create backup: %v", err), MB_OK|MB_ICONINFORMATION)
		}
		os.Exit(1)
	}
	fmt.Printf("Backup created at: %s\n", backupFile)

	// Read and parse the Preferences file
	preferences, err := readPreferencesFile(preferencesFile)
	if err != nil {
		fmt.Printf("Error reading Preferences file: %v\n", err)
		if runtime.GOOS == "windows" {
			if ShowMessageBox("Error", fmt.Sprintf("Failed to read Preferences: %v\nRestore from backup?", err), MB_OKCANCEL|MB_ICONQUESTION) == IDOK {
				restoreFromBackup(backupFile, preferencesFile)
			}
		} else {
			fmt.Print("Restore from backup? (y/n): ")
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			if strings.ToLower(strings.TrimSpace(response)) == "y" {
				restoreFromBackup(backupFile, preferencesFile)
			}
		}
		os.Exit(1)
	}

	// Ensure the required structure exists
	ensureStructureExists(preferences)

	// Generate a unique key for the model
	modelKey := fmt.Sprintf("custom:venice_%s", uuid.New().String()[:8])

	// Check if Venice.AI model already exists
	veniceModelIndex := -1
	for i, model := range preferences.Brave.AIChat.CustomModels {
		if model.Label == VeniceLabel {
			veniceModelIndex = i
			modelKey = model.Key // Keep the existing key
			break
		}
	}

	// Create the model configuration
	modelConfig := CustomModel{
		APIKey:          *apiKey,
		ContextSize:     VeniceContextSize,
		EndpointURL:     VeniceEndpoint,
		Key:             modelKey,
		Label:           VeniceLabel,
		ModelRequestName: VeniceModel,
	}

	// Update or add the Venice.AI model
	if veniceModelIndex >= 0 {
		preferences.Brave.AIChat.CustomModels[veniceModelIndex] = modelConfig
		fmt.Println("Updated existing Venice.AI configuration.")
	} else {
		preferences.Brave.AIChat.CustomModels = append(preferences.Brave.AIChat.CustomModels, modelConfig)
		fmt.Println("Added new Venice.AI configuration.")
	}

	// Set as default model
	preferences.Brave.AIChat.DefaultModelKey = modelKey

	// Write back to the Preferences file
	err = writePreferencesFile(preferencesFile, preferences)
	if err != nil {
		fmt.Printf("Error writing Preferences file: %v\n", err)
		if runtime.GOOS == "windows" {
			if ShowMessageBox("Error", fmt.Sprintf("Failed to write Preferences: %v\nRestore from backup?", err), MB_OKCANCEL|MB_ICONQUESTION) == IDOK {
				restoreFromBackup(backupFile, preferencesFile)
			}
		} else {
			fmt.Print("Restore from backup? (y/n): ")
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			if strings.ToLower(strings.TrimSpace(response)) == "y" {
				restoreFromBackup(backupFile, preferencesFile)
			}
		}
		os.Exit(1)
	}

	fmt.Println("Successfully updated Leo AI helper to use Venice.AI!")
	fmt.Println("IMPORTANT: If Brave is still running, you MUST completely close and restart it for changes to take effect.")
	if runtime.GOOS == "windows" {
		ShowMessageBox("Success - "+AppVersion, "Leo AI updated to use Venice.AI!\n\nIMPORTANT: If Brave is still running, you MUST completely close and restart it for changes to take effect.", MB_OK|MB_ICONINFORMATION)
	}
}