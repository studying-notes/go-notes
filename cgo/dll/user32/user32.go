package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	YesNoCancel = 0x00000003
)

func main() {
	fmt.Println("---------Start---------")

	ShowMessageA()
	ShowMessageB()
	ShowMessageC()
	ShowMessageD()

	fmt.Println("---------Stop---------")
}

// The first DLL method call
func ShowMessageA() int {
	user32, _ := syscall.LoadLibrary("user32.dll")
	msgBox, _ := syscall.GetProcAddress(user32, "MessageBoxW")
	defer syscall.FreeLibrary(user32)

	ret, _, err := syscall.Syscall9(msgBox, 4, 0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("CGO DLL"))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("The first DLL method call"))),
		YesNoCancel, 0, 0, 0, 0, 0,
	)
	if err != 0 {
		panic(err.Error())
	}
	return int(ret)
}

// The second DLL method call
func ShowMessageB() {
	user32 := syscall.NewLazyDLL("user32.dll")
	MessageBoxW := user32.NewProc("MessageBoxW")
	MessageBoxW.Call(
		uintptr(0),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("The second DLL method call"))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("CGO DLL"))),
		uintptr(0),
	)
}

// The third DLL method call
func ShowMessageC() {
	user32, _ := syscall.LoadDLL("user32.dll")
	MessageBoxW, _ := user32.FindProc("MessageBoxW")
	MessageBoxW.Call(
		uintptr(0),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("The third DLL method call"))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("CGO DLL"))),
		uintptr(0),
	)
}

// The fourth DLL method call
func ShowMessageD() {
	user32 := syscall.MustLoadDLL("user32.dll")
	MessageBoxW := user32.MustFindProc("MessageBoxW")
	MessageBoxW.Call(
		uintptr(0),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("The fourth DLL method call"))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("CGO DLL"))),
		uintptr(0),
	)
}
