package main

// #include <stdio.h>
// #include <malloc.h>
// typedef struct {
// char *msg;
// } Bug;
//
// void bug(Bug *b) {
// printf("%s", b->msg);
// }
import "C"
import (
	"runtime"
	"unsafe"
)

func main() {
	bug := C.Bug{C.CString("Hello, World!")}
	runtime.SetFinalizer(&bug, func(bug *C.Bug) {
		C.free(unsafe.Pointer(bug.msg))
	})
	C.bug(&bug)
	runtime.KeepAlive(&bug)
}
