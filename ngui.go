// Copyright (c) 2014 The cef2go authors. All rights reserved.
// License: BSD 3-clause.
// Website: https://github.com/CzarekTomczak/cef2go
// 			https://github.com/nvsoft/ngui

package main

import (
	"github.com/nvsoft/cef"
	"log"
	"os"
	"syscall"
	//"time"
	"unsafe"
	"github.com/nvsoft/wingui"
	"fmt"
)

var Logger *log.Logger = log.New(os.Stdout, "[main] ", log.Lshortfile)

func main() {
	hInstance, e := wingui.GetModuleHandle(nil)
	if e != nil {
		wingui.AbortErrNo("GetModuleHandle", e)
	}

	// you need to register to the callback before we fork processes
	cef.RegisterV8Callback("sup", cef.V8Callback(func(args []cef.V8Value) {
			arg0 := cef.V8ValueToInt32(args[0])
			arg1 := cef.V8ValueToInt32(args[1])
			arg2 := cef.V8ValueToBool(args[2])
			arg3 := cef.V8ValueToString(args[3])
			fmt.Printf("Calling V8Callback args: %d %d %v %s\n", arg0, arg1, arg2, arg3)
		}))

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

	wndproc := syscall.NewCallback(WndProc)
	Logger.Println("CreateWindow")
	hwnd := wingui.CreateWindow("ngui window", wndproc)

	browserSettings := cef.BrowserSettings{}
	// TODO: It should be executable's directory used
	// rather than working directory.
	url, _ := os.Getwd()
	url = "file://" + url + "/example.html"
	//cef.CreateBrowser(unsafe.Pointer(hwnd), &browserSettings, url, true)
	go func() {
		browser := cef.CreateBrowser(unsafe.Pointer(hwnd), &browserSettings, url, false)
		cef.WindowResized(unsafe.Pointer(hwnd))
		browser.ExecuteJavaScript("console.log('we outchea');cef2go.callback('sup', 10, 10, true, 'something');", "sup.js", 1)
	}()

	// It should be enough to call WindowResized after 10ms,
	// though to be sure let's extend it to 100ms.
	//time.AfterFunc(time.Millisecond*100, func() {
	//		cef.WindowResized(unsafe.Pointer(hwnd))
	//	})

	cef.RunMessageLoop()
	cef.Shutdown()
	os.Exit(0)
}

func WndProc(hwnd syscall.Handle, msg uint32, wparam, lparam uintptr) (rc uintptr) {
	switch msg {
	case wingui.WM_CREATE:
		rc = wingui.DefWindowProc(hwnd, msg, wparam, lparam)
	case wingui.WM_SIZE:
		// 最小化时不能调整Cef窗体，否则恢复时界面一片空白
		if (wparam == wingui.SIZE_RESTORED || wparam == wingui.SIZE_MAXIMIZED) {
			cef.WindowResized(unsafe.Pointer(hwnd))
		}
	case wingui.WM_CLOSE:
		wingui.DestroyWindow(hwnd)
	case wingui.WM_DESTROY:
		cef.QuitMessageLoop()
	default:
		rc = wingui.DefWindowProc(hwnd, msg, wparam, lparam)
	}
	return
}
