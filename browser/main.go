package main

import (
	"github.com/nvsoft/cef"
	"github.com/nvsoft/win"
	"unsafe"
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
