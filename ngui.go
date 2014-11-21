// Copyright (c) 2014 The ngui authors. All rights reserved.
// License: BSD 3-clause.
// Website: https://github.com/CzarekTomczak/cef2go
// 			https://github.com/nvsoft/ngui

package ngui

import (
	"errors"
	"fmt"
	"github.com/nvsoft/cef"
	"github.com/nvsoft/win"
	"log"
	"os"
	"strings"
	"syscall"
	//"time"
	"unsafe"
)

const (
	ICON_MAIN = 101

	SPI_SETDRAGFULLWINDOWS =      0x0025
)

var hInstance win.HINSTANCE
var Logger *log.Logger = log.New(os.Stdout, "[main] ", log.Lshortfile)
var wndProc = syscall.NewCallback(WndProc)
var transparentWndProc = syscall.NewCallback(TransparentWndProc)
var manifest Manifest

var rcWindow win.RECT
var dragPoint win.POINT
var isEnableDrag = true
var isDrag bool = false

const nguiWindowClass = `\o/ NGui_Window_Class \o/`
const nguiTransparentWindowClass = `\o/ NGui_Transparent_Window_Class \o/`

func init() {
	hInstance := win.GetModuleHandle(nil)
	if hInstance == 0 {
		panic("GetModuleHandle")
	}
	manifest.Load()
	MustRegisterWindowClass(nguiWindowClass)
	MustRegisterTransparentWindowClass(nguiTransparentWindowClass)
}

type Application struct {
}

func (this *Application) init() (err error) {
	cef.ExecuteProcess(unsafe.Pointer(hInstance))

	settings := cef.Settings{}
	settings.CachePath = manifest.CachePath()      // Set to empty to disable
	settings.LogSeverity = cef.LOGSEVERITY_DEFAULT // LOGSEVERITY_VERBOSE
	//settings.LocalesDirPath = releasePath + "/locales"
	settings.Locale = manifest.Locale() //"zh-CN"
	//settings.RemoteDebuggingPort = 7000
	cef.Initialize(settings)

	return
}

func MustRegisterWindowClass(className string) {
	//hIcon := win.LoadIcon(0, (*uint16)(unsafe.Pointer(uintptr(win.IDI_APPLICATION))))
	hIcon, _ := NewIconFromResource(hInstance, ICON_MAIN)
	if hIcon == 0 {
		panic("LoadIcon")
	}

	hCursor := win.LoadCursor(0, (*uint16)(unsafe.Pointer(uintptr(win.IDC_ARROW))))
	if hCursor == 0 {
		panic("LoadCursor")
	}

	var wc win.WNDCLASSEX
	wc.CbSize = uint32(unsafe.Sizeof(wc))
	wc.LpfnWndProc = wndProc
	wc.HInstance = hInstance
	wc.HIcon = hIcon
	wc.HCursor = hCursor
	wc.HbrBackground = win.COLOR_WINDOW + 1 //COLOR_BTNFACE
	wc.LpszClassName = syscall.StringToUTF16Ptr(className)

	if atom := win.RegisterClassEx(&wc); atom == 0 {
		panic("RegisterClassEx")
	}
}

func MustRegisterTransparentWindowClass(className string) {
	//hIcon := win.LoadIcon(0, (*uint16)(unsafe.Pointer(uintptr(win.IDI_APPLICATION))))
	hIcon, _ := NewIconFromResource(hInstance, ICON_MAIN)
	if hIcon == 0 {
		panic("LoadIcon")
	}

	hCursor := win.LoadCursor(0, (*uint16)(unsafe.Pointer(uintptr(win.IDC_ARROW))))
	if hCursor == 0 {
		panic("LoadCursor")
	}

	var wc win.WNDCLASSEX
	wc.CbSize = uint32(unsafe.Sizeof(wc))
	wc.LpfnWndProc = transparentWndProc
	wc.HInstance = hInstance
	wc.HIcon = hIcon
	wc.HCursor = hCursor
	wc.HbrBackground = win.BS_SOLID//win.COLOR_WINDOW + 1 //COLOR_BTNFACE
	wc.LpszClassName = syscall.StringToUTF16Ptr(className)

	if atom := win.RegisterClassEx(&wc); atom == 0 {
		panic("RegisterClassEx")
	}
}

// 创建浏览器窗口
func (this *Application) CreateBrowserWindow(url string, enable_transparent bool) (err error) {
	var dwExStyle, dwStyle uint32 = 0, 0
	fmt.Printf("CreateBrowserWindow enable_transparent=%v\n", enable_transparent)
	if enable_transparent {
		//dwExStyle = 0//win.WS_EX_LAYERED

		// 无边框效果不错。但是不能移动
		//dwStyle = win.WS_POPUP //& ^ (win.WS_CAPTION | win.WS_BORDER)

		dwStyle = win.WS_MINIMIZEBOX | win.WS_MAXIMIZEBOX | win.WS_CAPTION | win.WS_CLIPCHILDREN |
				win.WS_SYSMENU | win.WS_POPUP

	} else {
		dwStyle = win.WS_OVERLAPPEDWINDOW
	}

	if !manifest.FormFixed() {
		dwStyle |= win.WS_SIZEBOX
	}

	// 获取屏幕宽度和高度
	var x, y int32
	var width, height int32

	width = manifest.Width()
	height = manifest.Height()
	x = (win.GetSystemMetrics(win.SM_CXSCREEN) - width) / 2
	y = (win.GetSystemMetrics(win.SM_CYSCREEN) - height) / 2 - 2

	renderWindow := win.CreateWindowEx(
		dwExStyle,
		syscall.StringToUTF16Ptr(nguiTransparentWindowClass),
		nil,
		dwStyle, //|win.WS_CLIPSIBLINGS,
		x,       //win.CW_USEDEFAULT,
		y,       //win.CW_USEDEFAULT,
		width,   //win.CW_USEDEFAULT,
		height,  //win.CW_USEDEFAULT,
		0,       //hwndParent
		0,
		0, //hInstance
		nil)
	if renderWindow == 0 {
		err = errors.New("CreateWindowEx")
		return
	}

	if enable_transparent {
		/*
		SetWindowLong(hWnd, GWL_STYLE,
        GetWindowLong(hWnd, GWL_STYLE) & ~(WS_BORDER));
		 */

		// WS_BORDER Creates a window that has a border.创建一个有边界的窗口。
		// WS_EX_CLIENTEDGE Specifies that a window has a 3D look — that is, a border with a sunken edge.指定窗口具有3D外观，也即是一个下层的边缘。
		gwlStyle := win.GetWindowLong(renderWindow, win.GWL_STYLE)
		gwlStyle &= ^(win.WS_BORDER)
		//gwlStyle &= ^(win.WS_EX_CLIENTEDGE)
		win.SetWindowLong(renderWindow, win.GWL_STYLE, gwlStyle)

		//gwlExStyle := win.GetWindowLong(renderWindow, win.GWL_EXSTYLE)
		//gwlExStyle |= win.WS_EX_LAYERED;
		//win.SetWindowLong(renderWindow, win.GWL_EXSTYLE, gwlExStyle)

		// HWND hWnd: 窗口句柄
		// COLORREF col: 透明化颜色
		// byte bAlpha: ＝0：整个窗口透明, =255 完全不透明
		// uint32 dwFlags: ＝1：仅颜色 col 透明， =2 :窗口按照bAlpha变量进行透明处理。
		//win.SetLayeredWindowAttributes(renderWindow, 0, 255, 2)
		//win.UpdateWindow(renderWindow)
	}

	fmt.Printf("CreateBrowserWindow x=%v, y=%v, width=%v, height=%v\n", x, y, width, height)

	//win.ShowWindow(renderWindow, win.SW_SHOW)
	//win.UpdateWindow(renderWindow)

	browserSettings := cef.BrowserSettings{}

	go func() {

		//browser := cef.CreateBrowser(unsafe.Pointer(hwnd), &browserSettings, url, false)
		if enable_transparent {
			//win.MoveWindow(renderWindow, x, y, width, height, false)
			//win.SetWindowPos(renderWindow, 0, x, y, width, height, win.SWP_NOZORDER|win.SWP_NOACTIVATE|win.SWP_NOSIZE)
			win.SetWindowPos(renderWindow, 0, x, y, width, height, win.SWP_FRAMECHANGED)
		} else {
			win.MoveWindow(renderWindow, x, y, width, height, false)
		}

		cef.CreateBrowser(unsafe.Pointer(renderWindow), &browserSettings, url, false)

		//m_dwStyle = WS_CHILD | WS_CLIPCHILDREN | WS_CLIPSIBLINGS | WS_TABSTOP |
		//		WS_VISIBLE;
		cef.WindowResized(unsafe.Pointer(renderWindow))

		win.ShowWindow(renderWindow, win.SW_SHOW)//win.SW_SHOW
		//win.UpdateWindow(renderWindow)

		//win.SetWindowPos(renderWindow, 0, x, y, width, height, win.SWP_NOZORDER | win.SWP_NOACTIVATE | win.SWP_NOSIZE)

		//hIcon := win.LoadIcon(0, win.MAKEINTRESOURCE(win.IDI_ERROR))
		hIcon, _ := NewIconFromResource(hInstance, ICON_MAIN)
		if hIcon == 0 {
			panic("LoadIcon")
		}
		//win.SendMessage(renderWindow, win.WM_SETICON, 0, uintptr(hIcon))

		//cef.WindowResized(unsafe.Pointer(renderWindow))
		// It should be enough to call WindowResized after 10ms,
		// though to be sure let's extend it to 100ms.
		//time.AfterFunc(time.Millisecond*100, func() {
		//	cef.WindowResized(unsafe.Pointer(renderWindow))
		//})
	}()

	return
}

// 创建应用程序主窗口
func (this *Application) CreateBrowser() {
	url := manifest.FirstPage()
	wd, _ := os.Getwd()
	wd = strings.Replace(wd, "\\", "/", -1)
	url = "file:///" + wd + "/" + url
	fmt.Printf("CreateBrowser url=%s\n", url)
	enableTransparent := manifest.EnableTransparent()
	this.CreateBrowserWindow(url, enableTransparent)
}

func (e *Application) Exec() {
	cef.RunMessageLoop()
	cef.Shutdown()
	os.Exit(0)
}

func NewApplication() *Application {
	app := new(Application)
	app.init()

	return app
}

func WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	switch msg {
	case win.WM_CREATE:
		result = win.DefWindowProc(hwnd, msg, wParam, lParam)
	case win.WM_SIZE:
		// 最小化时不能调整Cef窗体，否则恢复时界面一片空白
		if wParam == win.SIZE_RESTORED || wParam == win.SIZE_MAXIMIZED {
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

func TransparentWndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	switch msg {
	case win.WM_CREATE:
		result = win.DefWindowProc(hwnd, msg, wParam, lParam)
		break
	case win.WM_NCPAINT:
		result = win.DefWindowProc(hwnd, msg, wParam, lParam)
		break
	case win.WM_LBUTTONDOWN:
		fmt.Printf("WM_LBUTTONDOWN\n")
		if isEnableDrag {
			win.SetCapture(hwnd)

			win.GetWindowRect(hwnd, &rcWindow)
			win.GetCursorPos(&dragPoint)

			isDrag = true
		}
		return 0
	case win.WM_LBUTTONUP:
		if win.GetCapture() == hwnd {
			win.ReleaseCapture()
		}
		isDrag = false
		return 0
	case win.WM_NCCALCSIZE:
		var size_param *win.NCCALCSIZE_PARAMS = (*win.NCCALCSIZE_PARAMS)(unsafe.Pointer(lParam));
		if win.BOOL(wParam) == win.TRUE {
			//var mi win.MONITORINFO = win.GetMonitorInformation(hwnd)
			size_param.Rgrc[2] = size_param.Rgrc[1]
			size_param.Rgrc[1] = size_param.Rgrc[0]
		}
		result = 0
		break
	case win.WM_NCHITTEST:
		//x := win.LOWORD(uint32(lParam))
		//y := win.HIWORD(uint32(lParam))
		//s := fmt.Sprintf("WM_NCHITTEST x,y=%v,%v\n", x, y)

		result = win.DefWindowProc(hwnd, msg, wParam, lParam)
		break
	case win.WM_MOUSEMOVE:
		//fmt.Printf("WM_MOUSEMOVE\n")
		if isDrag {
			var pe win.POINT
			win.GetCursorPos(&pe); // The new position for cursor pointer

			left := rcWindow.Left + (pe.X - dragPoint.X) // The horizontal position of the new window
			top := rcWindow.Top + (pe.Y - dragPoint.Y) // The vertical position of the new windows

			//win.MoveWindow(hwnd,reWindow.Left,reWindow.Top,reWindow.Right,reWindow.Bottom,true);// Moving window
			win.SetWindowPos(hwnd, 0, left, top, 0, 0, win.SWP_NOSIZE | win.SWP_NOZORDER)
		}
		break
	case win.WM_SIZE:
		// 最小化时不能调整Cef窗体，否则恢复时界面一片空白
		if wParam == win.SIZE_RESTORED || wParam == win.SIZE_MAXIMIZED {
			cef.WindowResized(unsafe.Pointer(hwnd))
		}
		break
	case win.WM_CLOSE:
		win.DestroyWindow(hwnd)
		break
	case win.WM_DESTROY:
		cef.QuitMessageLoop()
		break
	default:
		result = win.DefWindowProc(hwnd, msg, wParam, lParam)
		break
	}
	//result = win.DefWindowProc(hwnd, msg, wParam, lParam)
	return
}
