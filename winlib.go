package winlib

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
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
func GetModuleFileName(handle windows.Handle) string {
	n := uint32(1024)
	var buf []uint16
	for {
		buf = make([]uint16, n)
		r, err := windows.GetModuleFileName(handle, &buf[0], n)
		if err != nil {
			return ""
		}
		if r < n {
			break
		}
		// r == n means n not big enough
		n += 1024
	}
	return syscall.UTF16ToString(buf)
}

type json_client struct {
	_call uintptr /* Proc */
}

func (it *json_client) init(_dllName string) {
	global.Echo(_dllName, "_dllName")
	var handle windows.Handle
	var err error
	if filepath.IsAbs(_dllName) {
		global.Echo("<isAbs>")
		handle, err = windows.LoadLibraryEx(
			_dllName,
			0,
			windows.LOAD_WITH_ALTERED_SEARCH_PATH)
	} else {
		global.Echo("NOT <isAbs>")
		handle, err = windows.LoadLibrary(_dllName)
	}
	//global.Echo(handle, "handle")
	//global.Echo(err, "err")
	if err != nil {
		panic(err)
	}
	it._call, _ = windows.GetProcAddress(handle, "Call")
}

func NewJsonClient(_dllName string) *json_client {
	it := new(json_client)
	it.init(_dllName)
	return it
}

func (it *json_client) Call(name string, args any) (any, error) {
	_json := global.ToJson(args)
	ptr, _, _ := syscall.Syscall(it._call, 0,
		StringToUTF8Addr(name),
		StringToUTF8Addr(_json),
		0)
	output := UTF8AddrToString(ptr)
	output = strings.TrimSpace(output)
	if strings.HasPrefix(output, "\"") {
		result := global.FromJson(output)
		err_msg := result.(string)
		return nil, errors.New(err_msg)
	} else if strings.HasPrefix(output, "[") {
		result := global.FromJson(output)
		var ary []any = result.([]any)
		return ary[0], nil
	} else {
		panic(fmt.Sprintf("%s() returned malformed result json %s", name, output))
	}
}

type json_server struct {
	_funcTable map[string]func(any) (any, error)
}

func (it *json_server) init() {
	it._funcTable = make(map[string]func(any) (any, error))
}

func NewJsonServer() *json_server {
	it := new(json_server)
	it.init()
	return it
}

func (it *json_server) Register(_name string, _func func(any) (any, error)) {
	it._funcTable[_name] = _func
}

func (it *json_server) HandleCall(_namePtr, _jsonPtr uintptr) uintptr {
	_name := UTF8AddrToString(_namePtr)
	if it._funcTable[_name] == nil {
		err_msg := fmt.Sprintf("%s() not defined", _name)
		_output := global.ToJson(err_msg)
		return StringToUTF8Addr(_output)
	}
	_json := UTF8AddrToString(_jsonPtr)
	_input := global.FromJson(_json)
	_answer, _error := it._funcTable[_name](_input)
	if _error != nil {
		_output := global.ToJson(_error.Error())
		return StringToUTF8Addr(_output)
	} else {
		_output := global.ToPrettyJson(_answer)
		return StringToUTF8Addr(_output)
	}
}
