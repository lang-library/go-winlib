package winlib

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

func StringToWideCharPtr(s string) uintptr {
	arr, _ := windows.UTF16PtrFromString(s)
	return uintptr(unsafe.Pointer(arr))
}

func WideCharPtrToString(s uintptr) string {
	p := (*uint16)(unsafe.Pointer(s))
	return windows.UTF16PtrToString(p)
}