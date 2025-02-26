package main

import (
    "bufio"
    "bytes"
    "encoding/json"
    "flag"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "runtime"
    "strings"
    "time"

    "github.com/google/uuid"
)

// Windows API constants and types for MessageBox (Windows only)
// These are only used on Windows
var (
    MB_OK              uint32 = 0x00000000
    MB_OKCANCEL        uint32 = 0x00000001
    MB_YESNOCANCEL     uint32 = 0x00000003
    MB_ICONINFORMATION uint32 = 0x00000040
    MB_ICONQUESTION    uint32 = 0x00000020
    IDOK               int    = 1
    IDCANCEL           int    = 2
    IDYES              int    = 6
    IDNO               int    = 7
)

// Configuration settings
const (
    VeniceLabel       = "Venice.AI"
    VeniceModel       = "llama-3.3-70b"
    VeniceEndpoint    = "https://api.venice.ai/api/v1/chat/completions"
    VeniceContextSize = 4000
)

// CustomModel represents a Leo AI custom model configuration
type CustomModel struct {
    APIKey          string `json:"api_key"`
    ContextSize     int    `json:"context_size"`
    EndpointURL     string `json:"endpoint_url"`
    Key             string `json:"key"`
    Label           string `json:"label"`
    ModelRequestName string `json:"model_request_name"`
}

// BravePreferences represents the structure of Brave's Preferences.json file
type BravePreferences struct {
    Brave struct {
        AIChat struct {
            CustomModels   []CustomModel `json:"custom_models"`
            DefaultModelKey string       `json:"default_model_key"`
        } `json:"ai_chat"`
    } `json:"brave"`
}

// ShowMessageBox displays a Windows MessageBox (Windows only)
// This is a stub that will be replaced by the platform-specific implementation
func ShowMessageBox(title, text string, flags uint32) int {
    return 0 // Default implementation does nothing
}

// GetAPIKeyFromDialog shows a dialog (Windows) or CLI prompt (other platforms)
func GetAPIKeyFromDialog() (string, bool) {
    if runtime.GOOS == "windows" {
        // Windows-specific PowerShell dialog
        psScript := `
Add-Type -AssemblyName System.Windows.Forms
Add-Type -AssemblyName System.Drawing

$form = New-Object System.Windows.Forms.Form
$form.Text = "Venice.AI API Key Required"
$form.Size = New-Object System.Drawing.Size(500, 320)
$form.StartPosition = "CenterScreen"
$form.FormBorderStyle = "FixedDialog"
$form.MaximizeBox = $false
$form.MinimizeBox = $false

$label = New-Object System.Windows.Forms.Label
$label.Location = New-Object System.Drawing.Point(20, 20)
$label.Size = New-Object System.Drawing.Size(460, 20)
$label.Text = "Please enter your Venice.AI API key to configure Leo AI helper."
$form.Controls.Add($label)

$linkLabel = New-Object System.Windows.Forms.LinkLabel
$linkLabel.Location = New-Object System.Drawing.Point(20, 50)
$linkLabel.Size = New-Object System.Drawing.Size(460, 20)
$linkLabel.Text = "Click here to generate a new API key at venice.ai/settings/api"
$linkLabel.LinkArea = New-Object System.Windows.Forms.LinkArea(0, 10)
$linkLabel.add_LinkClicked({
    [System.Diagnostics.Process]::Start("https://venice.ai/settings/api")
})
$form.Controls.Add($linkLabel)

$inputLabel = New-Object System.Windows.Forms.Label
$inputLabel.Location = New-Object System.Drawing.Point(20, 90)
$inputLabel.Size = New-Object System.Drawing.Size(460, 20)
$inputLabel.Text = "API Key:"
$form.Controls.Add($inputLabel)

$textBox = New-Object System.Windows.Forms.TextBox
$textBox.Location = New-Object System.Drawing.Point(20, 120)
$textBox.Size = New-Object System.Drawing.Size(460, 20)
$form.Controls.Add($textBox)

$okButton = New-Object System.Windows.Forms.Button
$okButton.Location = New-Object System.Drawing.Point(310, 160)
$okButton.Size = New-Object System.Drawing.Size(80, 30)
$okButton.Text = "Submit"
$okButton.DialogResult = [System.Windows.Forms.DialogResult]::OK
$form.AcceptButton = $okButton
$form.Controls.Add($okButton)

$cancelButton = New-Object System.Windows.Forms.Button
$cancelButton.Location = New-Object System.Drawing.Point(400, 160)
$cancelButton.Size = New-Object System.Drawing.Size(80, 30)
$cancelButton.Text = "Cancel"
$cancelButton.DialogResult = [System.Windows.Forms.DialogResult]::Cancel
$form.CancelButton = $cancelButton
$form.Controls.Add($cancelButton)

$creditLabel = New-Object System.Windows.Forms.Label
$creditLabel.Location = New-Object System.Drawing.Point(20, 240)
$creditLabel.Size = New-Object System.Drawing.Size(460, 20)
$creditLabel.Text = "Created by George Larson - twitter.com/g3ologic - github.com/georgeglarson - george.g.larson@gmail.com"
$creditLabel.Font = New-Object System.Drawing.Font("Arial", 8)
$creditLabel.ForeColor = [System.Drawing.Color]::Gray
$form.Controls.Add($creditLabel)

$form.Topmost = $true
$form.Add_Shown({$textBox.Focus()})
$result = $form.ShowDialog()

if ($result -eq [System.Windows.Forms.DialogResult]::OK) {
    $textBox.Text
} else {
    ""
}
`
        tempDir := os.TempDir()
        psFile := filepath.Join(tempDir, "venice_api_dialog.ps1")
        err := os.WriteFile(psFile, []byte(psScript), 0644)
        if err != nil {
            fmt.Printf("Error creating PowerShell script: %v\n", err)
            return "", false
        }
        defer os.Remove(psFile)

        cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-WindowStyle", "Hidden", "-File", psFile)
        var out bytes.Buffer
        var stderr bytes.Buffer
        cmd.Stdout = &out
        cmd.Stderr = &stderr
        err = cmd.Run()
        if err != nil {
            fmt.Printf("Error executing PowerShell script: %v\nSTDERR: %s\n", err, stderr.String())
            return "", false
        }

        result := strings.TrimSpace(out.String())
        if result == "" {
            return "", false
        }
        return result, true
    }

    // CLI prompt for macOS and Linux
    fmt.Println("Please enter your Venice.AI API key (or press Enter to cancel):")
    fmt.Println("You can find or generate your API key at: https://venice.ai/settings/api")
    reader := bufio.NewReader(os.Stdin)
    input, err := reader.ReadString('\n')
    if err != nil {
        fmt.Printf("Error reading input: %v\n", err)
        return "", false
    }
    result := strings.TrimSpace(input)
    if result == "" {
        return "", false
    }
    return result, true
}

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

    fmt.Println("Leo Venice.AI Configuration Tool")
    fmt.Println("================================")

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
        ShowMessageBox("Success", "Leo AI updated to use Venice.AI!\n\nIMPORTANT: If Brave is still running, you MUST completely close and restart it for changes to take effect.", MB_OK|MB_ICONINFORMATION)
    }
}

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