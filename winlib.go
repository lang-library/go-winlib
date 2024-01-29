package winlib

import (
	"syscall"
	"unsafe"

	"github.com/lang-library/go-global"
	"golang.org/x/sys/windows"
)

func StringToWideAddr(s string) uintptr {
	arr := StringToWidePtr(s)
	return uintptr(unsafe.Pointer(arr))
}

func WideAddrToString(s uintptr) string {
	p := (*uint16)(unsafe.Pointer(s))
	return WidePtrToString(p)
}

func StringToWidePtr(s string) *uint16 {
	p, err := windows.UTF16PtrFromString(s)
	if err != nil {
		return nil
	}
	return p
}

func WidePtrToString(p *uint16) string {
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

type json_api struct {
	_call *syscall.Proc
}

func (it *json_api) init(_dllName string) {
	_dll, _ := syscall.LoadDLL(_dllName)
	it._call, _ = _dll.FindProc("Call")
}

func NewJsonAPI(_dllName string) *json_api {
	self := new(json_api)
	self.init(_dllName)
	return self
}

func (it *json_api) Call(name string, args any) any {
	_json := global.ToJson(args)
	ptr, _, _ := it._call.Call(
		StringToUTF8Addr(name),
		StringToUTF8Addr(_json))
	output := UTF8AddrToString(ptr)
	result := global.FromJson(output)
	return result
}
