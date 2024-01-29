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
	it := new(json_api)
	it.init(_dllName)
	return it
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

func (it *json_api) CallOne(name string, args any) any {
	_result := it.Call(name, args)
	if _result == nil {
		return nil
	}
	var _ary []any
	var _ok bool
	if _ary, _ok = _result.([]any); !_ok {
		return _result
	}
	return _ary[0]
}

type json_server struct {
	_funcTable map[string]func(any) any
}

func (it *json_server) init() {
	it._funcTable = make(map[string]func(any) any)
}

func NewJsonServer() *json_server {
	it := new(json_server)
	it.init()
	return it
}

func (it *json_server) Register(_name string, _func func(any) any) {
	it._funcTable[_name] = _func
}

func (it *json_server) HandleCall(_namePtr, _jsonPtr uintptr) uintptr {
	_name := UTF8AddrToString(_namePtr)
	if it._funcTable[_name] == nil {
		return StringToUTF8Addr("null")
	}
	_json := UTF8AddrToString(_jsonPtr)
	_input := global.FromJson(_json)
	_answer := it._funcTable[_name](_input)
	_output := global.ToPrettyJson(_answer)
	return StringToUTF8Addr(_output)
}
