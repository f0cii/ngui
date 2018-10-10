package main

import (
	"unsafe"

	"github.com/sumorf/cef"
	"github.com/sumorf/win"
)

var (
	hInstance win.HINSTANCE
)

func init() {
	hInstance := win.GetModuleHandle(nil)
	if hInstance == 0 {
		panic("GetModuleHandle")
	}
}

func main() {
	cef.ExecuteProcess(unsafe.Pointer(hInstance))
}
