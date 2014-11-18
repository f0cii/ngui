package ngui

import (
	"errors"
	"fmt"
	"github.com/nvsoft/win"
	"path/filepath"
	"strings"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

var (
	kernel                = syscall.MustLoadDLL("kernel32.dll")
	getModuleFileNameProc = kernel.MustFindProc("GetModuleFileNameW")
)

func ExePath() string {
	exePath, _ := Executable()
	exeDir := filepath.Dir(exePath)
	return exeDir
}

// GetModuleFileName() with hModule = NULL
func Executable() (exePath string, err error) {
	return getModuleFileName()
}

func getModuleFileName() (string, error) {
	var n uint32
	b := make([]uint16, syscall.MAX_PATH)
	size := uint32(len(b))

	r0, _, e1 := getModuleFileNameProc.Call(0, uintptr(unsafe.Pointer(&b[0])), uintptr(size))
	n = uint32(r0)
	if n == 0 {
		return "", e1
	}
	s := string(utf16.Decode(b[0:n]))
	s = strings.Replace(s, "\\", "/", -1)
	return s, nil
}

func NewIconFromResource(instance win.HINSTANCE, resId uint16) (ico win.HICON, err error) {
	if ico = win.LoadIcon(instance, win.MAKEINTRESOURCE(uintptr(resId))); ico == 0 {
		err = errors.New(fmt.Sprintf("Cannot load icon from resource with id %v", resId))
	}

	return ico, err
}
