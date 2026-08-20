//go:build windows

package main

import (
	"github.com/wailsapp/wails/v3/pkg/application"
	"golang.org/x/sys/windows"
)

var (
	user32           = windows.NewLazySystemDLL("user32.dll")
	setWindowPosProc = user32.NewProc("SetWindowPos")
)

func setWindowAlwaysOnBottom(window *application.WebviewWindow) {
	if window == nil {
		return
	}
	hwnd := uintptr(window.NativeWindow())
	if hwnd == 0 {
		return
	}
	const (
		hwndBottom    = 1
		swpNoSize     = 0x0001
		swpNoMove     = 0x0002
		swpNoActivate = 0x0010
	)
	setWindowPosProc.Call(hwnd, hwndBottom, 0, 0, 0, 0, swpNoSize|swpNoMove|swpNoActivate)
}
