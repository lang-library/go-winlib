package winlib

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

func StringToUTF16Addr(s string) uintptr {
	arr := StringToUTF16Ptr(s)
	return uintptr(unsafe.Pointer(arr))
}

func UTF16AddrToString(s uintptr) string {
	p := (*uint16)(unsafe.Pointer(s))
	return UTF16PtrToString(p)
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

func StringToUTF8Addr(s string) uintptr {
	arr := StringToUTF8Ptr(s)
	return uintptr(unsafe.Pointer(arr))
}

func UTF8AddrToString(s uintptr) string {
	p := (*byte)(unsafe.Pointer(s))
	return UTF8PtrToString(p)
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
		if char == 0 {
			break
		}
		chars = append(chars, char)
	}
	return string(chars)
}

// https://go.dev/src/os/executable_windows.go
func GetModuleFileName(handle windows.Handle) string /*, error*/ {
	n := uint32(1024)
	var buf []uint16
	for {
		buf = make([]uint16, n)
		r, err := windows.GetModuleFileName(handle, &buf[0], n)
		if err != nil {
			return "" /*, err*/
		}
		if r < n {
			break
		}
		// r == n means n not big enough
		n += 1024
	}
	return syscall.UTF16ToString(buf) /*, nil*/
}
