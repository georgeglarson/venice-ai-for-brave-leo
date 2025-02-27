package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// findPreferencesFile finds Brave's Preferences.json file based on the OS
func findPreferencesFile() (string, error) {
	var possiblePaths []string
	switch runtime.GOOS {
	case "windows":
		localAppData := os.Getenv("LOCALAPPDATA")
		possiblePaths = []string{
			filepath.Join(localAppData, "BraveSoftware", "Brave-Browser", "User Data", "Default", "Preferences"),
			filepath.Join(localAppData, "Brave Software", "Brave-Browser", "User Data", "Default", "Preferences"),
		}
	case "darwin": // macOS
		homeDir, _ := os.UserHomeDir()
		possiblePaths = []string{
			filepath.Join(homeDir, "Library", "Application Support", "BraveSoftware", "Brave-Browser", "Default", "Preferences"),
			filepath.Join(homeDir, "Library", "Application Support", "Brave Software", "Brave-Browser", "Default", "Preferences"),
		}
	default: // Linux and others
		homeDir, _ := os.UserHomeDir()
		possiblePaths = []string{
			filepath.Join(homeDir, ".config", "BraveSoftware", "Brave-Browser", "Default", "Preferences"),
			filepath.Join(homeDir, ".config", "brave-browser", "Default", "Preferences"),
		}
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not find Brave browser's Preferences file")
}

// isBraveRunning checks if Brave browser is currently running (including background processes)
func isBraveRunning() bool {
	var runningProcesses []string
	
	switch runtime.GOOS {
	case "windows":
		possibleProcesses := []string{
			"brave.exe",
			"brave-browser.exe",
			"BraveBrowser.exe",
			"Brave.exe",
			"Brave-Browser.exe",
		}
		for _, process := range possibleProcesses {
			cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("IMAGENAME eq %s", process), "/NH")
			output, err := cmd.Output()
			if err == nil && len(output) > 0 && !strings.Contains(string(output), "No tasks") && !strings.Contains(string(output), "INFO: No tasks") {
				runningProcesses = append(runningProcesses, process)
			}
		}
	case "darwin": // macOS
		cmd := exec.Command("pgrep", "-i", "brave")
		output, err := cmd.Output()
		if err == nil && len(output) > 0 {
			runningProcesses = append(runningProcesses, "brave")
		}
	default: // Linux and others
		cmd := exec.Command("pgrep", "-i", "brave")
		output, err := cmd.Output()
		if err == nil && len(output) > 0 {
			runningProcesses = append(runningProcesses, "brave")
		}
	}
	
	return len(runningProcesses) > 0
}

// killBraveProcesses attempts to kill all Brave browser processes
// Returns true if successful, false otherwise
func killBraveProcesses() bool {
	success := true
	
	switch runtime.GOOS {
	case "windows":
		possibleProcesses := []string{
			"brave.exe",
			"brave-browser.exe",
			"BraveBrowser.exe",
			"Brave.exe",
			"Brave-Browser.exe",
		}
		for _, process := range possibleProcesses {
			cmd := exec.Command("taskkill", "/F", "/IM", process)
			if err := cmd.Run(); err != nil {
				// Ignore errors as some processes might not exist
				fmt.Printf("Note: Could not kill %s (may not be running)\n", process)
			} else {
				fmt.Printf("Killed process: %s\n", process)
			}
		}
	case "darwin": // macOS
		cmd := exec.Command("pkill", "-9", "-i", "brave")
		if err := cmd.Run(); err != nil {
			fmt.Println("Note: Could not kill Brave processes (may not be running)")
			success = false
		} else {
			fmt.Println("Killed Brave processes")
		}
	default: // Linux and others
		cmd := exec.Command("pkill", "-9", "-i", "brave")
		if err := cmd.Run(); err != nil {
			fmt.Println("Note: Could not kill Brave processes (may not be running)")
			success = false
		} else {
			fmt.Println("Killed Brave processes")
		}
	}
	
	// Verify all processes were killed
	time.Sleep(500 * time.Millisecond) // Give OS time to update process list
	return !isBraveRunning() && success
}