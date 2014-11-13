package ngui

import (
	"fmt"
	"unsafe"
	"github.com/nvsoft/win"
	"github.com/nvsoft/cef"
)

var printf = fmt.Println

func registerV8Handlers() {
	cef.RegisterV8Handler("move", v8_move)
}

func v8_move(browser *cef.Browser, args []cef.V8Value) {
	fmt.Println("v8_move")
	x := cef.V8ValueToInt32(args[0])
	y := cef.V8ValueToInt32(args[1])
	fmt.Printf("v8_move x=%v,y=%v\n", x, y)
	hWnd := browser.GetWindowHandle()

	h := (win.HWND)(unsafe.Pointer(hWnd))
	var rect win.RECT
	win.GetWindowRect(h, &rect)
	width := int32(rect.Right - rect.Left)
	height := int32(rect.Bottom - rect.Top)

	win.MoveWindow(h, x, y, width, height, false)
	//cef.WindowResized(unsafe.Pointer(hWnd))
}
