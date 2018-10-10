package ngui

import (
	"fmt"
	"os"

	"github.com/sumorf/cef"
	"github.com/sumorf/win"
)

var printf = fmt.Println

const (
	SC_DRAGMOVE = 0xF012
)

// 调用Js
//browser.ExecuteJavaScript("console.log('something from go invoke');alert('something from go invoke');", "go.js", 1)

// 初始化Handler
func init() {
	cef.RegisterV8Handler("createWindow", win_handler_create_window)
	cef.RegisterV8Handler("startDrag", win_handler_startDrag)
	cef.RegisterV8Handler("restore", win_handler_restore)
	cef.RegisterV8Handler("minimize", win_handler_minimize)
	cef.RegisterV8Handler("maximize", win_handler_maximize)
	cef.RegisterV8Handler("close", win_handler_close)
	cef.RegisterV8Handler("sizeTo", win_handler_sizeTo)     // width, hight
	cef.RegisterV8Handler("moveTo", win_handler_moveTo)     // left, top
	cef.RegisterV8Handler("setTitle", win_handler_setTitle) // title
	cef.RegisterV8Handler("quit", win_handler_quit)
}

// 创建窗口
func win_handler_create_window(browser *cef.Browser, args []cef.V8Value) (result interface{}) {
	fmt.Println("win_handler_create_window")

	url := cef.V8ValueToString(args[0])
	captionless := cef.V8ValueToBool(args[1])
	gApplication.CreateBrowserWindow(url, captionless)

	return
}

// 移动窗口
func win_handler_startDrag(browser *cef.Browser, args []cef.V8Value) (result interface{}) {
	fmt.Println("win_handler_startDrag")

	h := win.HWND(browser.GetWindowHandle())

	var pt win.POINT
	win.GetCursorPos(&pt)

	isDrag = true

	win.PostMessage(h, win.WM_LBUTTONDOWN, win.HTCAPTION, uintptr(win.MAKELONG(uint16(pt.X), uint16(pt.Y))))

	return
}

// 恢复窗口
func win_handler_restore(browser *cef.Browser, args []cef.V8Value) (result interface{}) {
	fmt.Println("win_handler_restore")

	h := win.HWND(browser.GetWindowHandle())
	win.ShowWindow(h, win.SW_RESTORE)

	return
}

// 最小化窗口
func win_handler_minimize(browser *cef.Browser, args []cef.V8Value) (result interface{}) {
	fmt.Println("win_handler_minimize")

	h := win.HWND(browser.GetWindowHandle())
	win.ShowWindow(h, win.SW_MINIMIZE)

	return
}

// 最大化窗口
func win_handler_maximize(browser *cef.Browser, args []cef.V8Value) (result interface{}) {
	fmt.Println("win_handler_maximize")

	h := win.HWND(browser.GetWindowHandle())
	win.ShowWindow(h, win.SW_MAXIMIZE)

	return
}

// 关闭窗口
func win_handler_close(browser *cef.Browser, args []cef.V8Value) (result interface{}) {
	h := win.HWND(browser.GetWindowHandle())
	win.SendMessage(h, win.WM_CLOSE, 0, 0)

	return
}

// 为窗口设置新的尺寸
func win_handler_sizeTo(browser *cef.Browser, args []cef.V8Value) (result interface{}) {
	fmt.Println("win_handler_sizeTo")
	width := cef.V8ValueToInt32(args[0])
	height := cef.V8ValueToInt32(args[1])

	h := win.HWND(browser.GetWindowHandle())
	var rect win.RECT
	win.GetWindowRect(h, &rect)

	fmt.Printf("win_handler_sizeTo Left=%v,Right=%v,Width=%v,Height=%v\n", rect.Left, rect.Top, width, height)
	win.MoveWindow(h, rect.Left, rect.Top, width, height, true)

	//result = 1

	return
}

// 为窗口设置新的位置
func win_handler_moveTo(browser *cef.Browser, args []cef.V8Value) (result interface{}) {
	fmt.Println("win_handler_moveTo")
	left := cef.V8ValueToInt32(args[0])
	top := cef.V8ValueToInt32(args[1])

	fmt.Printf("win_handler_moveTo left=%v,top=%v\n", left, top)

	h := win.HWND(browser.GetWindowHandle())

	var rect win.RECT
	win.GetWindowRect(h, &rect)
	width := int32(rect.Right - rect.Left)
	height := int32(rect.Bottom - rect.Top)

	fmt.Printf("win_handler_moveTo Left=%v,Right=%v,Width=%v,Height=%v\n", left, top, width, height)
	win.MoveWindow(h, left, top, width, height, true)

	return
}

// 为窗口设置标题
func win_handler_setTitle(browser *cef.Browser, args []cef.V8Value) (result interface{}) {
	title := cef.V8ValueToString(args[0])

	h := win.HWND(browser.GetWindowHandle())
	win.SetWindowText(h, title)

	return
}

// 退出程序
func win_handler_quit(browser *cef.Browser, args []cef.V8Value) (result interface{}) {
	h := win.HWND(browser.GetWindowHandle())
	win.SendMessage(h, win.WM_CLOSE, 0, 0)
	//win.PostQuitMessage(0);
	os.Exit(1)

	return
}
