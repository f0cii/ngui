// Copyright (c) 2014 The ngui authors. All rights reserved.
// License: BSD 3-clause.
// Website: https://github.com/CzarekTomczak/cef2go
// 			https://github.com/nvsoft/ngui

package ngui

import (
	"errors"
	"log"
	"os"
	"syscall"
	"time"
	"unsafe"
	"github.com/nvsoft/cef"
	"github.com/nvsoft/win"
)

var IDR_MAINFRAME = win.MAKEINTRESOURCE(100)
var Logger *log.Logger = log.New(os.Stdout, "[main] ", log.Lshortfile)
var wndproc = syscall.NewCallback(WndProc)

const nguiWindowClass = `\o/ NGui_Window_Class \o/`

func init() {
	MustRegisterWindowClass(nguiWindowClass)
}

type Engine struct {
}

func (this *Engine) init() (err error) {
	hInstance := win.GetModuleHandle(nil)
	if hInstance == 0 {
		err = errors.New("GetModuleHandle")
		return
	}

	// you need to register to the callback before we fork processes
	/*cef.RegisterV8Callback("sup", cef.V8Callback(func(args []cef.V8Value) {
			arg0 := cef.V8ValueToInt32(args[0])
			arg1 := cef.V8ValueToInt32(args[1])
			arg2 := cef.V8ValueToBool(args[2])
			arg3 := cef.V8ValueToString(args[3])
			fmt.Printf("Calling V8Callback args: %d %d %v %s\n", arg0, arg1, arg2, arg3)
		}))*/

	cef.ExecuteProcess(unsafe.Pointer(hInstance))

	settings := cef.Settings{}
	settings.CachePath = "webcache"                // Set to empty to disable
	settings.LogSeverity = cef.LOGSEVERITY_DEFAULT // LOGSEVERITY_VERBOSE
	//settings.ResourcesDirPath = releasePath
	//settings.LocalesDirPath = releasePath + "/locales"
	//settings.CachePath = cwd + "/webcache"      // Set to empty to disable
	//settings.LogSeverity = cef.LOGSEVERITY_INFO // LOGSEVERITY_VERBOSE
	//settings.LogFile = cwd + "/debug.log"
	//settings.RemoteDebuggingPort = 7000
	cef.Initialize(settings)

	return
}

func MustRegisterWindowClass(className string) {
	hInst := win.GetModuleHandle(nil)
	if hInst == 0 {
		panic("GetModuleHandle")
	}

	hIcon := win.LoadIcon(0, (*uint16)(unsafe.Pointer(uintptr(win.IDI_APPLICATION))))
	if hIcon == 0 {
		panic("LoadIcon")
	}

	hCursor := win.LoadCursor(0, (*uint16)(unsafe.Pointer(uintptr(win.IDC_ARROW))))
	if hCursor == 0 {
		panic("LoadCursor")
	}

	var wc win.WNDCLASSEX
	wc.CbSize = uint32(unsafe.Sizeof(wc))
	wc.LpfnWndProc = wndproc
	wc.HInstance = hInst
	wc.HIcon = hIcon
	wc.HCursor = hCursor
	wc.HbrBackground = win.COLOR_BTNFACE+1
	wc.LpszClassName = syscall.StringToUTF16Ptr(className)

	if atom := win.RegisterClassEx(&wc); atom == 0 {
		panic("RegisterClassEx")
	}
}

func CreateWindowEx(title string, wndproc uintptr) (hwnd syscall.Handle, err error) {
	var hwndParent win.HWND = 0

	hWnd := win.CreateWindowEx(
		0,
		syscall.StringToUTF16Ptr(nguiWindowClass),
		nil,
		win.WS_OVERLAPPEDWINDOW, //|win.WS_CLIPSIBLINGS,
		win.CW_USEDEFAULT,
		win.CW_USEDEFAULT,
		win.CW_USEDEFAULT,
		win.CW_USEDEFAULT,
		hwndParent,
		0,
		0,
		nil)
	if hWnd == 0 {
		err = errors.New("CreateWindowEx")
		return
	}

	// ShowWindow
	win.ShowWindow(hWnd, win.SW_SHOWDEFAULT)

	hwnd = syscall.Handle(hWnd)

	return
}

func (this *Engine) CreateWindow(url string) {
	Logger.Println("CreateWindow")
	hwnd, _ := CreateWindowEx("ngui window", wndproc)
	//wh := win.HWND(hwnd)
	//win.ShowWindow(wh, win.SW_HIDE)

	browserSettings := cef.BrowserSettings{}

	//cef.CreateBrowser(unsafe.Pointer(hwnd), &browserSettings, url, true)
	go func() {
		//browser := cef.CreateBrowser(unsafe.Pointer(hwnd), &browserSettings, url, false)
		cef.CreateBrowser(unsafe.Pointer(hwnd), &browserSettings, url, false)
		//cef.WindowResized(unsafe.Pointer(hwnd))
		//win.ShowWindow(wh, win.SW_SHOWNORMAL)
		//browser.ExecuteJavaScript("console.log('we outchea');ngui.callback('sup', 10, 10, true, 'something');", "sup.js", 1)
	}()

	// It should be enough to call WindowResized after 10ms,
	// though to be sure let's extend it to 100ms.
	time.AfterFunc(time.Millisecond*100, func() {
			cef.WindowResized(unsafe.Pointer(hwnd))
		})
}

func (e *Engine) Exec() {
	cef.RunMessageLoop()
	cef.Shutdown()
	os.Exit(0)
}

func NewEngine() *Engine {
	e := new(Engine)

	e.init()

	return e
}

func WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	switch msg {
	case win.WM_CREATE:
		result = win.DefWindowProc(hwnd, msg, wParam, lParam)
	case win.WM_SIZE:
		// 最小化时不能调整Cef窗体，否则恢复时界面一片空白
		if (wParam == win.SIZE_RESTORED || wParam == win.SIZE_MAXIMIZED) {
			cef.WindowResized(unsafe.Pointer(hwnd))
		}
	case win.WM_CLOSE:
		win.DestroyWindow(hwnd)
	case win.WM_DESTROY:
		cef.QuitMessageLoop()
	default:
		result = win.DefWindowProc(hwnd, msg, wParam, lParam)
	}
	return
}
