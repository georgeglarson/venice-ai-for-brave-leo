//go:build windows
// +build windows

package main

import (
	"syscall"
	"unsafe"
)

// ShowMessageBox displays a Windows MessageBox
// Windows-specific implementation
func ShowMessageBox(title, text string, flags uint32) int {
	user32 := syscall.NewLazyDLL("user32.dll")
	getActiveWindow := user32.NewProc("GetActiveWindow")
	messageBox := user32.NewProc("MessageBoxW")
	hwnd, _, _ := getActiveWindow.Call()
	ret, _, _ := messageBox.Call(
		hwnd,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(text))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))),
		uintptr(flags),
	)
	return int(ret)
}