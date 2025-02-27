//go:build !windows
// +build !windows

package main

// ShowMessageBox displays a message box (stub implementation for non-Windows platforms)
func ShowMessageBox(title, text string, flags uint32) int {
	return 0 // Default implementation does nothing
}