package main

import (
	"fmt"
	"reflect"
)

func main() {
	ta := reflect.ArrayOf(10, reflect.TypeOf(0))
	tc := reflect.ChanOf(reflect.BothDir, reflect.TypeOf(0))
	tp := reflect.PtrTo(reflect.TypeOf(0))
	ts := reflect.SliceOf(reflect.TypeOf(0))
	tm := reflect.MapOf(reflect.TypeOf(0), reflect.TypeOf(0))
	tf := reflect.FuncOf([]reflect.Type{reflect.TypeOf(0)}, []reflect.Type{reflect.TypeOf(0)}, false)
	tt := reflect.StructOf([]reflect.StructField{
		{Name: "A", Type: reflect.TypeOf(0)},
		{Name: "B", Type: reflect.TypeOf(0)},
	})

	reflect.MakeChan(tc, 5)

	fmt.Println(ta, tc, tp, ts, tm, tf, tt)
}
