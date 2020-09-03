package pkg

// #cgo CFLAGS: -I./include
// #include <stdlib.h>
// #include "base64.h"
import "C"
import "unsafe"

func Base64Encode(input string) string {
	cInput := C.CString(input)
	cOutput := C.base64_encode(cInput)
	defer func() {
		C.free(unsafe.Pointer(cInput)) // <stdlib.h>
		C.free(unsafe.Pointer(cOutput))
	}()
	return C.GoString(cOutput)
}

func Base64Decode(input string) string {
	cInput := C.CString(input)
	cOutput := C.base64_decode(cInput)
	defer func() {
		C.free(unsafe.Pointer(cInput))
		C.free(unsafe.Pointer(cOutput))
	}()
	return C.GoString(cOutput)
}
