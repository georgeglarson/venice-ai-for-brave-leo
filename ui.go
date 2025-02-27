package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// GetAPIKeyFromDialog shows a dialog (Windows) or CLI prompt (other platforms)
func GetAPIKeyFromDialog() (string, bool) {
	if runtime.GOOS == "windows" {
		// Windows-specific PowerShell dialog
		psScript := `
Add-Type -AssemblyName System.Windows.Forms
Add-Type -AssemblyName System.Drawing

$form = New-Object System.Windows.Forms.Form
$form.Text = "Venice.AI API Key Required - ${AppVersion}"
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
$creditLabel.Location = New-Object System.Drawing.Point(20, 220)
$creditLabel.Size = New-Object System.Drawing.Size(460, 40)
$creditLabel.Text = "Created by George Larson (github.com/georgeglarson)
@g3ologic"
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