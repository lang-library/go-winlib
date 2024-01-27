package winlib

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

func StringToWideCharAddr(s string) uintptr {
	arr, _ := windows.UTF16PtrFromString(s)
	return uintptr(unsafe.Pointer(arr))
}

func WideCharAddrToString(s uintptr) string {
	p := (*uint16)(unsafe.Pointer(s))
	return windows.UTF16PtrToString(p)
}

func StringToUTF16Ptr(s string) *uint16 {
	p, err := windows.UTF16PtrFromString(s)
	if err != nil {
		return nil
	}
	return p
}

func UTF16PtrToString(p *uint16) string {
	return windows.UTF16PtrToString(p)
}

func StringToUTF8Ptr(s string) *byte {
	p, err := syscall.BytePtrFromString(s)
	if err != nil {
		return nil
	}
	return p
}

func UTF8PtrToString(p *byte) string {
	if p == nil {
		return ""
	}
	var char byte
	var chars = []byte{}
	for i := 0; ; i++ {
		char = *(*byte)(unsafe.Pointer(unsafe.Add(unsafe.Pointer(p), unsafe.Sizeof(byte(0))*uintptr(i))))
		// null char
		if char == 0 {
			break
		}
		chars = append(chars, char)
	}
	return string(chars)
}
