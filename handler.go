package ngui

import (
	"fmt"
	//"unsafe"
	//"github.com/nvsoft/win"
	"github.com/nvsoft/cef"
)

var printf = fmt.Println

func init() {
	cef.RegisterV8Handler("move", v8_move)	// x, y
	cef.RegisterV8Handler("close", v8_close)
}

func v8_move(browser *cef.Browser, args []cef.V8Value) (result interface {}) {
	fmt.Println("v8_move")
	x := cef.V8ValueToInt32(args[0])
	y := cef.V8ValueToInt32(args[1])

	fmt.Printf("v8_move x=%v,y=%v\n", x, y)

	// 调用Js
	browser.ExecuteJavaScript("console.log('something from go invoke');alert('something from go invoke');", "go.js", 1)

	/*hWnd := browser.GetWindowHandle()

	h := (win.HWND)(hWnd)
	var rect win.RECT
	win.GetWindowRect(h, &rect)
	width := int32(rect.Right - rect.Left)
	height := int32(rect.Bottom - rect.Top)

	win.MoveWindow(h, x, y, width, height, false)
	*/

	//cef.WindowResized(unsafe.Pointer(hWnd))
	result = 3

	return
}

func v8_close(browser *cef.Browser, args []cef.V8Value) (result interface {}) {

	return
}
